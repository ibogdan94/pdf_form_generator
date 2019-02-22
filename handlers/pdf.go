package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pdf_form_generator/parser"
	"strings"
)

type Handlers struct {
	pdfParser *parser.PdfParser
	baseUrl   string
}

func (h Handlers) upload(ctx *gin.Context) {
	file, headers, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "PDF file not found",
		})
		return
	}

	mimeType := headers.Header.Get("Content-Type")

	if mimeType != "application/pdf" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Only PDP file can be loaded",
		})
		return
	}

	defer file.Close()

	pages, err := h.pdfParser.PdfToPng(file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	var result []string

	for _, page := range pages {
		result = append(result, h.convertToPublicUrl(page))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (h Handlers) show(ctx *gin.Context) {
	storeFolderName := ctx.Param("folderName")

	var pageElements parser.PngToPdf

	if err := ctx.BindJSON(&pageElements); err != nil {
		log.Printf("Json error:%s\n", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot bind pdf params",
		})
		return
	}

	result, err := h.pdfParser.PngsToPdf(storeFolderName, pageElements)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": h.convertToPublicUrl(result),
	})
}

func (h Handlers) convertToPublicUrl(absolutePath string) string {
	return strings.Replace(absolutePath, "/go/src/pdf_form_generator", h.baseUrl, -1)
}

func (h *Handlers) SetupRoutes(r *gin.Engine) {
	r.POST("/api/pdf", h.upload)
	r.GET("/api/pdf/:folderName", h.show)
}

func NewHandlers(pwd string, baseUrl string) *Handlers {
	return &Handlers{
		&parser.PdfParser{
			pwd,
			pwd + "/" + "store",
			pwd + "/" + "result",
			nil,
		},
		baseUrl,
	}
}
