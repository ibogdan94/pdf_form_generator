package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pdf_form_generator/parser"
)

type Handlers struct {
	pdfParser *parser.PdfParser
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

	result, err := h.pdfParser.PdfToPng(file)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (h Handlers) show(ctx *gin.Context) {
	storeFolderName := ctx.Param("folderName")

	result, err := h.pdfParser.PngsToPdf(storeFolderName)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (h *Handlers) SetupRoutes(r *gin.Engine) {
	r.POST("/pp", h.upload)
	r.GET("/pp/:folderName", h.show)
}

func NewHandlers(pwd string) *Handlers {
	return &Handlers{
		&parser.PdfParser{
			pwd,
			pwd + "/" + "store",
			pwd + "/" + "result",
		},
	}
}
