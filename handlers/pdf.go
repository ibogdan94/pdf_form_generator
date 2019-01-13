package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pdf_form_generator/parser"
)

type Handlers struct {
	pwd string
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

	p := parser.PdfParser{
		h.pwd,
		h.pwd + "/" + "output",
	}

	result, err := p.PdfToPng(file)

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (h *Handlers) SetupRoutes(r *gin.Engine) {
	r.POST("/pp", h.upload)
}

func NewHandlers(pwd string) *Handlers {
	return &Handlers{
		pwd,
	}
}
