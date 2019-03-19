package parser

import (
	"context"
	"errors"
	"fmt"
	"github.com/unidoc/unidoc/pdf/creator"
	"golang.org/x/sync/errgroup"
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
	Code         *string
}

func (p *PdfParser) setCode(code string) {
	p.Code = &code
}

func (p PdfParser) getCode() string {
	return *p.Code
}

func (p PdfParser) PdfToPng(file multipart.File) (result []string, err error) {
	p.setCode(utils.Random())
	folderToSave := p.StoreFolder + "/" + p.getCode()

	if err := p.createFolder(folderToSave); err != nil {
		return result, errors.New(fmt.Sprintf("Cannot create result folder %s. Error: %s", folderToSave, err))
	}

	absolutePath := folderToSave + "/" + p.getCode() + ".pdf"
	f, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		return result, errors.New(fmt.Sprintf("Cannot Open File. Error: %s", err))
	}

	defer f.Close()
	io.Copy(f, file)

	numberOfPages, err := GetNumberOfPages(absolutePath)

	if err != nil {
		return result, errors.New(fmt.Sprintf("Cannot get number of pages. Error: %s", err))
	}

	if numberOfPages == 0 {
		return result, errors.New("pdf file has no pages")
	}

	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < int(numberOfPages); i++ {
		page := fmt.Sprintf("%s[%d]", absolutePath, i)
		output := fmt.Sprintf("%s/%s[%d].png", folderToSave, p.getCode(), i)
		//300dpi for printing
		command := fmt.Sprintf("convert -resize 2480x3508 -verbose -trim -density %d -depth %d -flatten %s %s", 300, 8, page, output)
		result = append(result, output)

		g.Go(func() error {
			return p.converterCommand(command)
		})
	}

	if err := g.Wait(); err != nil {
		return result, err
	}

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

func (p PdfParser) converterCommand(cmd string) (err error) {
	parts := strings.Fields(cmd)
	args := parts[1:]
	_, err = exec.Command(parts[0], args...).Output()

	return err
}

func (p PdfParser) createFolder(outputFolderPath string) (err error) {
	if _, err := os.Stat(outputFolderPath); os.IsNotExist(err) {
		if err := os.Mkdir(outputFolderPath, 0755); err != nil {
			return err
		}
	}
	return
}

func (p PdfParser) PngsToPdf(code string, pageElements PngToPdf) (result string, err error) {
	storeFolder := p.StoreFolder + "/" + code

	if _, err := os.Stat(storeFolder); os.IsNotExist(err) {
		log.Printf("storeFolder error: %s", err)
		return result, err
	}

	numberOfPages := len(pageElements.Pages)

	if numberOfPages == 0 {
		return result, errors.New("no pages")
	}

	p.setCode(code)

	resultFolder := p.ResultFolder + "/" + p.getCode()

	if _, err := os.Stat(resultFolder); !os.IsNotExist(err) {
		if err := os.RemoveAll(resultFolder); err != nil {
			return result, errors.New("cannot remove folder")
		}

	}

	if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
		if err := os.Mkdir(resultFolder, 0755); err != nil {
			return result, err
		}
	}

	pngs := make(map[int]string, numberOfPages)

	g, _ := errgroup.WithContext(context.Background())

	for index, page := range pageElements.Pages {
		g.Go(func() error {
			resultPage, err := p.addPlaceholdersToPngImage(page, pageElements.Data)

			if err != nil {
				return err
			}

			pngs[index] = resultPage

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return result, err
	}

	result, err = p.generatePdf(sortPages(pngs))

	if err != nil {
		return result, err
	}

	return result, nil
}

func (p PdfParser) addPlaceholdersToPngImage(pngPage PngPageWithElements, data []DataType) (fileDestination string, err error) {
	var args []string

	args = append(args, fmt.Sprintf("%s/%s/%s[%d].png", p.StoreFolder, p.getCode(), p.getCode(), pngPage.Page-1))

	if len(data) == 0 && len(pngPage.CanvasElements.Objects) == 0 {
		return fileDestination, errors.New("not enough data to render page")
	}

	wg := new(sync.WaitGroup)
	c := make(chan []string)

	wg.Add(len(pngPage.CanvasElements.Objects))

	usedTokens := make([]string, len(data))

	//@todo need to normalize after the frontend will be updated
	for _, object := range pngPage.CanvasElements.Objects {
		go func(c chan<- []string, data []DataType, object Text) {
			for _, dataType := range data {
				if stringInSlice(dataType.Token, usedTokens) {
					continue
				}

				if dataType.Value == object.Text && dataType.Placeholder == "" {
					usedTokens = append(usedTokens, dataType.Token)

					c <- []string{
						"-fill",
						object.Fill,
						"-undercolor",
						object.BackgroundColor,
						"-pointsize",
						fmt.Sprintf("%v", object.FontSize),
						"-weight",
						fmt.Sprintf("%v", object.FontWeight),
						"-annotate",
						fmt.Sprintf("+%v+%v", object.Left, object.Top),
						object.Text,
					}

					break
				} else if dataType.Placeholder == object.Text {
					object.Text = dataType.Value
					usedTokens = append(usedTokens, dataType.Token)

					c <- []string{
						"-fill",
						object.Fill,
						"-undercolor",
						object.BackgroundColor,
						"-pointsize",
						fmt.Sprintf("%v", object.FontSize),
						"-weight",
						fmt.Sprintf("%v", object.FontWeight),
						"-annotate",
						fmt.Sprintf("+%v+%v", object.Left, object.Top),
						object.Text,
					}

					break
				}
			}
		}(c, data, object)
	}

	go func(c <-chan []string) {
		for argsFromChan := range c {
			args = append(args, argsFromChan...)
			wg.Done()
		}
	}(c)

	wg.Wait()
	close(c)

	fileDestination = fmt.Sprintf("%s/%s/%s[%d].png", p.ResultFolder, p.getCode(), p.getCode(), pngPage.Page-1)
	args = append(args, []string{"-verbose", "-trim", "-density", "300", "-depth", "8", "-flatten", fileDestination}...)

	cmd := exec.Command("magick", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fileDestination, errors.New(string(output))
	}

	return fileDestination, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

func (p PdfParser) generatePdf(pngsAbsolutePath []string) (pdfFile string, err error) {
	c := creator.New()

	for _, imgPath := range pngsAbsolutePath {
		img, err := creator.NewImageFromFile(imgPath)

		if err != nil {
			return pdfFile, err
		}

		img.ScaleToWidth(612.0)

		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		err = c.Draw(img)

		if err != nil {
			return pdfFile, err
		}
	}

	pdfFile = p.ResultFolder + "/" + p.getCode() + "/" + p.getCode() + ".pdf"

	if err := c.WriteToFile(pdfFile); err != nil {
		return pdfFile, err
	}

	return pdfFile, nil
}
