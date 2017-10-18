package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	utils "../utils"
	"os"
	"github.com/gorilla/sessions"
	"bytes"
	"image/png"
	"encoding/base64"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type B64 struct {
	B64 string `form:"b64" json:"b64" binding:"required"`
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

	result, err := utils.ParsePdfToPng(file, headers)

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
	var json B64

	if err := ctx.BindJSON(&json); err != nil {
		fmt.Println("Json error:", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot bind JSON",
		})
	}

	var b64FromRequest string = json.B64
	//remove extra js headers
	var b64 string = b64FromRequest[22:]

	unBased, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		fmt.Println("Cannot decode b64. Error:", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot decode b64",
		})
	}

	r := bytes.NewReader(unBased)
	img, err := png.Decode(r)

	if err != nil {
		fmt.Println("Bad png:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad png",
		})
	}

	folderForResult := "./temp/result"

	if stat, err := os.Stat(folderForResult); err != nil && stat.IsDir() == false {
		if err := os.Mkdir(folderForResult, 0755); err != nil {
			fmt.Println("Cannot Create Folder:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
		}
	}

	fileName := utils.Random()
	targetPath := folderForResult + "/" + fileName + ".png"

	f, err := os.Create(targetPath)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	if err := f.Close(); err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	pdfPath := folderForResult + "/" + fileName + ".pdf"

	inputPaths := []string{targetPath}

	err = utils.ImagesToPdf(inputPaths, pdfPath)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"url": pdfPath[1:],
	})
}