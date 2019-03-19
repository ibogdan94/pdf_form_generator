package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pdf_form_generator/parser"
	"strings"
)

type Handlers struct {
	pdfParser *parser.PdfParser
	baseUrl   string
}

func (h Handlers) upload(ctx *gin.Context) {
	file, headers, err := ctx.Request.FormFile("pdf")

	if hasErrorHandler(ctx, err, "PDF file not found", http.StatusBadRequest) {
		return
	}

	mimeType := headers.Header.Get("Content-Type")

	if mimeType != "application/pdf" {
		hasErrorHandler(ctx, err, "Only PDP file can be loaded", http.StatusBadRequest)
		return
	}

	defer file.Close()

	pages, err := h.pdfParser.PdfToPng(file)

	if hasErrorHandler(ctx, err, fmt.Sprintf("Something went wrong. %s", err), http.StatusInternalServerError) {
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
		hasErrorHandler(ctx, err, "Cannot bind pdf params", http.StatusBadRequest)
		return
	}

	result, err := h.pdfParser.PngsToPdf(storeFolderName, pageElements)

	if hasErrorHandler(ctx, err, "Cannot find the pdf file", http.StatusBadRequest) {
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
