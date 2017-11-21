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
	"regexp"
)

func ParsePdfToPng(schemaAndHost string, file multipart.File, code string, headers *multipart.FileHeader) (result []string, error error) {
	pwd, err := os.Getwd()

	if err != nil {
		return result, err
	}

	absoluteFolderForPDF := pwd + "/" + Config.TempPath + "/" + code
	relativeFolderForPDF := "/" + Config.TempPath + "/" + code

	if err := os.Mkdir(absoluteFolderForPDF, 0755); err != nil {
		log.Printf("Cannot Create Folder. Error: %s", err)
		return result, err
	}

	imageFilePathWithoutExtension := absoluteFolderForPDF + "/" + code
	extension := filepath.Ext(headers.Filename)
	imageFilePathWithExtension := imageFilePathWithoutExtension + extension
	f, err := os.OpenFile(imageFilePathWithExtension, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Printf("Cannot Open File. Error: %s", err)
		return result, err
	}

	defer f.Close()
	io.Copy(f, file)

	numberOfPages, err := GetNumberOfPages(imageFilePathWithExtension)

	if err != nil {
		return result, err
	}

	if numberOfPages > 0 {
		wg := new(sync.WaitGroup)

		for i := 0; i < int(numberOfPages); i++ {
			page := fmt.Sprintf("%s[%d]", imageFilePathWithExtension, i)
			output := fmt.Sprintf("%s/%s[%d].png", absoluteFolderForPDF, code, i)
			relativeOutput := fmt.Sprintf("%s/%s[%d].png", relativeFolderForPDF, code, i)
			//300dpi for printing
			command := fmt.Sprintf("convert -resize 2480x3508 -verbose -trim -density %d -depth %d -flatten %s %s", 300, 8, page, output)
			fmt.Println(command)
			wg.Add(1)
			result = append(result, schemaAndHost+relativeOutput)

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

func GetNumberOfPages(imageFilePathWithExtension string) (numberOfPages int, err error) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	if err := mw.PingImage(imageFilePathWithExtension); err != nil {
		fmt.Println("Cannot PingImage:", err)
		return numberOfPages, err
	}

	numberOfPages = int(mw.GetNumberImages())
	return
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

func ImagesWithPlaceHoldersToPdf(absoluteFilePath string, pngPage *PngPageWithElements, data []DataType) (fileDestination string, err error) {
	var args []string

	args = append(args, absoluteFilePath)

	if len(data) > 0 {
		for _, object := range pngPage.CanvasElements.Objects {
			//wg := new(sync.WaitGroup)

			for _, dataType := range data {
				//wg.Add(1)
				//go func(label *DataType) {
				if dataType.Placeholder == object.Text {
					object.Text = dataType.Value
					object.Prepare()

					args = append(args, "-fill", object.Fill, "-undercolor", object.BackgroundColor, "-pointsize", fmt.Sprintf("%v", object.FontSize), "-weight", fmt.Sprintf("%v", object.FontWeight), "-annotate", fmt.Sprintf("+%v+%v", object.Left, object.Top), object.Text)
				}
				//wg.Done()
				//}(&dataType)
			}

			//wg.Wait()
		}
	}

	fileDestination = absoluteFilePath[:len(absoluteFilePath)-4] + "_temp.png"
	args = append(args, []string{"-verbose", "-trim", "-density", "300", "-depth", "8", "-flatten", fileDestination,}...)

	cmd := exec.Command("magick", args...)

	fmt.Println(cmd)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf(fmt.Sprint(err) + ": " + string(output))
		return fileDestination, err
	}

	return fileDestination, err
}

func GetFilesInFolderByExt(absoluteDir string, ext string) (files []string, err error) {
	if _, err := os.Stat(absoluteDir); os.IsNotExist(err) {
		return files, err
	}

	filepath.Walk(absoluteDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})

	return files, err
}