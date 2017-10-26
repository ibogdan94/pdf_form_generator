package utils

import (
	"io"
	"fmt"
	"sync"
	"strings"
	"os/exec"
	"mime/multipart"
	"os"
	"path/filepath"
	"github.com/unidoc/unidoc/pdf/creator"
	"log"
	"gopkg.in/gographics/imagick.v3/imagick"
	"net/url"
)

func ParsePdfToPng(url *url.URL, file multipart.File, headers *multipart.FileHeader) (result []string, error error) {
	randomFileName := Random()

	props, err := ParseJSONConfig()

	if err != nil {
		log.Fatal(err)
	}

	folderForPDF := props.TempPath + randomFileName

	if err := os.Mkdir(folderForPDF, 0755); err != nil {
		fmt.Println("Cannot Create Folder:", err)
		return result, err
	}

	imageFilePathWithoutExtension := folderForPDF + "/" + randomFileName
	extension := filepath.Ext(headers.Filename)
	imageFilePathWithExtension := imageFilePathWithoutExtension + extension
	f, err := os.OpenFile(imageFilePathWithExtension, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("Cannot Open File:", err)
		return result, err
	}

	defer f.Close()
	io.Copy(f, file)

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	if err := mw.PingImage(imageFilePathWithExtension); err != nil {
		fmt.Println("Cannot PingImage:", err)
		return result, err
	}

	numberOfPages := mw.GetNumberImages()

	schemaAndHost := url.Scheme + "://" + url.Host

	if numberOfPages > 0 {
		wg := new(sync.WaitGroup)

		for i := 0; i < int(numberOfPages); i++ {
			page := fmt.Sprintf("%s[%d]", imageFilePathWithExtension, i)
			output := fmt.Sprintf("%s/%s.%d.png", folderForPDF, Random(), i)
			command := fmt.Sprintf("convert -verbose -trim -density %d -resize %s -depth %d -flatten %s %s", 400, "25%", 8, page, output)
			fmt.Println(command)
			wg.Add(1)
			result = append(result, schemaAndHost + output[1:])

			go pdfToImagesCommand(command, wg)
		}

		wg.Wait()
	}

	return result, err
}

func pdfToImagesCommand(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	args := parts[1:]
	out, err := exec.Command(parts[0], args...).Output()

	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}

	fmt.Printf("%s", out)
	wg.Done()
}

func ImagesToPdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for _, imgPath := range inputPaths {
		log.Printf("Image: %s", imgPath)

		img, err := creator.NewImageFromFile(imgPath)

		if err != nil {
			log.Printf("Error loading image: %v", err)
			return err
		}

		img.ScaleToWidth(612.0)

		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		err = c.Draw(img)

		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}
	}

	err := c.WriteToFile(outputPath)
	return err
}
