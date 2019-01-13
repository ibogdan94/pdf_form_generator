package parser

import (
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"pdf_form_generator/utils"
	"strings"
	"sync"
)

type Parser interface {
	PdfToPng(file multipart.File) ([]string, error)
	PngToPdf()
}

type PdfParser struct {
	Pwd          string
	OutputFolder string
}

func (p PdfParser) PdfToPng(file multipart.File) (result []string, error error) {
	filename := utils.Random()
	absolutePath := p.OutputFolder + "/" + filename + ".pdf"
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
		output := fmt.Sprintf("%s/%s[%d].png", p.OutputFolder, filename, i)
		//300dpi for printing
		command := fmt.Sprintf("convert -resize 2480x3508 -verbose -trim -density %d -depth %d -flatten %s %s", 300, 8, page, output)

		fmt.Println(command)

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
		fmt.Println("error")
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
		fmt.Println("error occured")
		fmt.Printf("%s", err)
		panic(err)
	}

	fmt.Printf("%s", out)
	wg.Done()
}

func (p PdfParser) PngToPdf() {

}
