package handlres

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"sync"
	"mime/multipart"
	utils "../utils"
	"os"
	"path/filepath"
	"io"
	"strings"
	"os/exec"
	"github.com/gorilla/sessions"
)

var SessionStore *sessions.FilesystemStore

func ValidateUploadPDF(c *gin.Context) {
	//file, headers, err := c.Request.FormFile("pdf")
	file, headers, err := c.Request.FormFile("pdf")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "PDF file not found",
		})
	}

	mimeType := headers.Header.Get("Content-Type")

	if mimeType != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Only PDP file can be loaded",
		})
	}

	defer file.Close()

	result, err := parsePDFToPng(file, headers)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Try again later",
		})
	}

	if len(result) > 0 {
		//fmt.Println("Inside len")
		//gfSession, err := SessionStore.Get(c.Request, "generated_imgs")
		//
		//if err != nil {
		//	fmt.Println("Session error:", err)
		//}
		//
		//gfSession.Values["images"] = result
		//
		//err = gfSession.Save(c.Request, c.Writer)
		//
		//fmt.Println("Save")
		//if err != nil {
		//	fmt.Println(err)
		//	c.Redirect(301, "/")
		//}
		//
		//fmt.Println("Redirect to edit")

		c.Redirect(301, "/pdf/edit")
	}

	c.Redirect(301, "/")
}

func EditParsedPDF(c *gin.Context) {
	gfSession, err := SessionStore.Get(c.Request, "generated_imgs")

	//no one images
	if err != nil {
		c.Redirect(301, "/")
	}

	c.HTML(http.StatusOK, "builder.html", map[string]interface{}{
		"images": gfSession,
	})
}

func parsePDFToPng(file multipart.File, headers *multipart.FileHeader) (result []string, error error) {
	var pngResult []string
	randomFileName := utils.Random()
	folderForPDF := "./temp/" + randomFileName

	if err := os.Mkdir(folderForPDF, 0755); err != nil {
		fmt.Println("Cannot Create Folder:", err)
		return pngResult, err
	}

	imageFilePathWithoutExtension := folderForPDF + "/" + randomFileName
	extension := filepath.Ext(headers.Filename)
	imageFilePathWithExtension := imageFilePathWithoutExtension+extension
	f, err := os.OpenFile(imageFilePathWithExtension, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println("Cannot Open File:", err)
		return pngResult, err
	}

	defer f.Close()
	io.Copy(f, file)

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	if err := mw.PingImage(imageFilePathWithExtension); err != nil {
		fmt.Println("Cannot PingImage:", err)
		return pngResult, err
	}

	numberOfPages := mw.GetNumberImages()

	fmt.Println("numberOfPages", numberOfPages)

	if numberOfPages > 0 {
		fmt.Println("Create wait group")
		wg := new(sync.WaitGroup)

		for i := 0; i < int(numberOfPages); i++ {
			page := fmt.Sprintf("%s[%d]", imageFilePathWithExtension, i)
			output := fmt.Sprintf("%s/%s.%d.png", folderForPDF, utils.Random(), i)
			fmt.Println("------------------")
			fmt.Println(output)
			fmt.Println("------------------")
			command := fmt.Sprintf("convert -density %d -resize %s -depth %d -flatten %s %s", 400, "25%", 8, page, output)
			fmt.Println(command)
			wg.Add(1)
			pngResult = append(pngResult, output)
			go exe_cmd(command, wg)
		}

		wg.Wait()
	}

	fmt.Println("PNG RESULT")
	fmt.Println(pngResult)

	return pngResult, err
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	args := parts[1:len(parts)]
	out, err := exec.Command(parts[0], args...).Output()
	if err != nil {
		fmt.Println("error occured")
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done()
}
