package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/location"
	"net/http"
	"fmt"
	"pdf_form_generator/utils"
	"os"
	"bytes"
	"image/png"
	"encoding/base64"
	"log"
)

type PngPage struct {
	B64 string `json:"b64" binding:"required"`
	Page int `json:"page" binding:"required"`
}

func ValidateUploadPDF(c *gin.Context) {
	file, headers, err := c.Request.FormFile("file")

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

	url := location.Get(c)

	result, err := utils.ParsePdfToPng(url, file, headers)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Try again later",
		})
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No results was generated",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"images": result,
	})
}

func SavePDF(ctx *gin.Context) {
	var json []PngPage

	if err := ctx.BindJSON(&json); err != nil {
		fmt.Println("Json error:", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot bind JSON",
		})
	}

	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Error:", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	props, err := utils.ParseJSONConfig()

	if err != nil {
		log.Fatal(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Config error",
		})
	}

	prefixName := utils.Random()

	resultFolder := pwd + "/" + props.TempPath + "/result"

	if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
		if err := os.Mkdir(resultFolder, 0755); err != nil {
			fmt.Println("Cannot Create Folder:", resultFolder, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
		}
	}

	folderForResult := resultFolder + "/" + prefixName
	relativeFolder := props.TempPath + "/result/" + prefixName

	if err := os.Mkdir(folderForResult, 0755); err != nil {
		fmt.Println("Cannot Create Folder:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	//array of generated png from request
	var inputPaths []string

	url := location.Get(ctx)

	schemaAndHost := url.Scheme + "://" + url.Host

	for _, pngPage := range json {
		var b64FromRequest string = pngPage.B64
		//remove extra js headers
		var b64 string = b64FromRequest[22:]

		unBased, err := base64.StdEncoding.DecodeString(b64)

		if err != nil {
			fmt.Println("Cannot decode b64. Error:", err)
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"message": "Cannot decode b64",
			//})
		}

		r := bytes.NewReader(unBased)
		img, err := png.Decode(r)

		if err != nil {
			fmt.Println("Bad png:", err)
			//ctx.JSON(http.StatusBadRequest, gin.H{
			//	"message": "Bad png",
			//})
		}

		targetPath := fmt.Sprintf("%s/%s[%d].png", folderForResult, prefixName, pngPage.Page)

		f, err := os.Create(targetPath)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			//ctx.JSON(http.StatusInternalServerError, gin.H{
			//	"message": "Something went wrong",
			//})
		}

		if err := png.Encode(f, img); err != nil {
			f.Close()
			fmt.Printf("Error: %v\n", err)
			//ctx.JSON(http.StatusInternalServerError, gin.H{
			//	"message": "Something went wrong",
			//})
		}

		if err := f.Close(); err != nil {
			fmt.Printf("Error: %v\n", err)
			//ctx.JSON(http.StatusInternalServerError, gin.H{
			//	"message": "Something went wrong",
			//})
		}

		relativePath := fmt.Sprintf("%s/%s[%d].png", relativeFolder, prefixName, pngPage.Page)

		inputPaths = append(inputPaths, relativePath)
	}

	pdfPath := fmt.Sprintf("%s/%s.pdf", folderForResult, prefixName)

	err = utils.ImagesToPdf(inputPaths, pdfPath)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	pdfRelativePath := fmt.Sprintf("%s/%s/%s.pdf", schemaAndHost, relativeFolder, prefixName)

	for i, image := range inputPaths {
		inputPaths[i] = schemaAndHost + "/" + image
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"pdf": pdfRelativePath,
		"images": inputPaths,
	})
}