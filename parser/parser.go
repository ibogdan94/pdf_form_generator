package parser

import (
	"fmt"
	"github.com/unidoc/unidoc/pdf/creator"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"pdf_form_generator/utils"
	"sort"
	"strings"
	"sync"
)

type Parser interface {
	PdfToPng(file multipart.File) ([]string, error)
	PngToPdf()
}

type PdfParser struct {
	Pwd          string
	StoreFolder  string
	ResultFolder string
}

func (p PdfParser) PdfToPng(file multipart.File) (result []string, err error) {
	name := utils.Random()
	folderToSave := p.StoreFolder + "/" + name

	if err := p.createFolder(folderToSave); err != nil {
		log.Println("Cannot create result folder:", folderToSave, err)
		return result, err
	}

	absolutePath := folderToSave + "/" + name + ".pdf"
	f, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		log.Printf("Cannot Open File. Error: %s", err)
		panic(err)
	}

	defer f.Close()
	io.Copy(f, file)

	numberOfPages, err := GetNumberOfPages(absolutePath)

	if err != nil {
		log.Printf("Cannot get number of pages. Error: %s", err)
		panic(err)
	}

	if numberOfPages == 0 {
		log.Printf("0 pages")
		panic("0 pages")
	}

	wg := new(sync.WaitGroup)

	for i := 0; i < int(numberOfPages); i++ {
		page := fmt.Sprintf("%s[%d]", absolutePath, i)
		output := fmt.Sprintf("%s/%s[%d].png", folderToSave, name, i)
		//300dpi for printing
		command := fmt.Sprintf("convert -resize 2480x3508 -verbose -trim -density %d -depth %d -flatten %s %s", 300, 8, page, output)

		wg.Add(1)
		result = append(result, output)

		go p.converterCommand(command, wg)
	}

	wg.Wait()

	return
}

func GetNumberOfPages(absolutePath string) (numberOfPages int, err error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	if err := mw.PingImage(absolutePath); err != nil {
		return numberOfPages, err
	}

	numberOfPages = int(mw.GetNumberImages())
	return
}

func (p PdfParser) converterCommand(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	args := parts[1:]
	out, err := exec.Command(parts[0], args...).Output()

	if err != nil {
		log.Printf("converterCommand error: %s", err)
		panic(err)
	}

	log.Printf("converterCommand output: %s", out)
	wg.Done()
}

func (p PdfParser) createFolder(outputFolderPath string) (err error) {
	if _, err := os.Stat(outputFolderPath); os.IsNotExist(err) {
		if err := os.Mkdir(outputFolderPath, 0755); err != nil {
			return err
		}
	}
	return
}

func (p PdfParser) PngsToPdf(code string) (result string, err error) {
	storeFolder := p.StoreFolder + "/" + code

	if _, err := os.Stat(storeFolder); os.IsNotExist(err) {
		log.Printf("storeFolder error: %s", err)
		return result, err
	}

	//todo code duplication
	resultFolder := p.ResultFolder + "/" + code

	if _, err := os.Stat(resultFolder); !os.IsNotExist(err) {
		return resultFolder + "/" + code + ".pdf", nil
	}

	pages := []int{0, 1, 2, 3}

	wg := new(sync.WaitGroup)
	wg.Add(len(pages))

	mutex := &sync.Mutex{}

	pngs := make(map[int]string, len(pages))

	for index, page := range pages {
		go func(index, page int, wg *sync.WaitGroup) {
			mutex.Lock()
			pngs[index] = fmt.Sprintf("%s/%s[%d].png", storeFolder, code, page)
			wg.Done()
			mutex.Unlock()
		}(index, page, wg)
	}

	wg.Wait()

	sortedPngs := sortPages(pngs)

	result, err = p.generatePdf(code, sortedPngs)

	if err != nil {
		return result, err
	}

	return result, nil
}

func sortPages(pngs map[int]string) []string {
	var keys []int

	for k := range pngs {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var sortedPngs []string

	for _, k := range keys {
		sortedPngs = append(sortedPngs, pngs[k])
	}

	return sortedPngs
}

func (p PdfParser) generatePdf(code string, pngsAbsolutePath []string) (pdfFile string, err error) {
	//todo code duplication
	resultFolder := p.ResultFolder + "/" + code

	if _, err := os.Stat(resultFolder); os.IsExist(err) {
		return resultFolder + "/" + code + ".pdf", nil
	}

	if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
		if err := os.Mkdir(resultFolder, 0755); err != nil {
			return pdfFile, err
		}
	}

	c := creator.New()

	for _, imgPath := range pngsAbsolutePath {
		//fmt.Printf("Image: %s", imgPath)

		img, err := creator.NewImageFromFile(imgPath)

		if err != nil {
			log.Printf("Error loading image: %v", err)
			return pdfFile, err
		}

		img.ScaleToWidth(612.0)

		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		err = c.Draw(img)

		if err != nil {
			log.Printf("Error: %v", err)
			return pdfFile, err
		}
	}

	pdfFile = resultFolder + "/" + code + ".pdf"

	if err := c.WriteToFile(pdfFile); err != nil {
		log.Printf("WriteToFile: %v", err)
		return pdfFile, err
	}

	return pdfFile, nil
}
