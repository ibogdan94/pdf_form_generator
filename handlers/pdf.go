package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"github.com/ibogdan94/pdf_form_generator/utils"
	"os"
	"log"
)

func ValidateUploadPDF(ctx *gin.Context) {
	file, headers, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "PDF file not found",
		})
	}

	mimeType := headers.Header.Get("Content-Type")

	if mimeType != "application/pdf" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Only PDP file can be loaded",
		})
	}

	defer file.Close()

	code := utils.Random()
	result, err := utils.ParsePdfToPng(utils.GetSchemaAndHost(ctx), file, code, headers)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Try again later",
		})
	}

	if len(result) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No results was generated",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"images": result,
		"code":   code,
	})
}

func GeneratePDF(ctx *gin.Context) {
	var pngData utils.PngToPdf
	code := ctx.Param("code")

	if err := ctx.BindJSON(&pngData); err != nil {
		log.Printf("Json error:%s\n", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot bind JSON",
		})
		return
	}

	pwd, err := os.Getwd()

	if err != nil {
		log.Printf("Error: %s\n", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	if len(pngData.Pages) == 0 {
		log.Printf("Error: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Empty data",
		})
		return
	}

	relativePrefix := utils.Config.TempPath + "/" + code + "/" + code

	//array of generated png from request
	var pngs []string

	for index, page := range pngData.Pages {
		pageName := fmt.Sprintf("%s[%d].png", relativePrefix, index)
		absolutePagePath := pwd + "/" + pageName
		if _, err := os.Stat(absolutePagePath); os.IsNotExist(err) {
			log.Printf("Page %d for file with code %s was not found. Error: %s\n", index, code, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Page does not exist",
			})
			return
		}

		fileDestination, err := utils.ImagesWithPlaceHoldersToPdf(absolutePagePath, &page, pngData.Data)

		if err != nil {
			log.Printf("Error: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		pngs = append(pngs, fileDestination)
	}

	pdfRelativeLink, err := savePDF(code, pngs)

	clearTempPngWithPlaceholders(pngs)

	if err != nil {
		log.Printf("Cannot generate pdf file: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot generate pdf file",
		})
		return
	}

	pdfUrl := utils.GetSchemaAndHost(ctx) + "/" + pdfRelativeLink

	ctx.JSON(http.StatusOK, gin.H{
		"pdf": pdfUrl,
	})
}

func clearTempPngWithPlaceholders(pngs []string) (err error) {
	for _, png := range pngs {
		err := os.Remove(png)
		if err != nil {
			return err
		}
	}

	return err
}

func savePDF(code string, pngsAbsolutePath []string) (pdfRelativeLink string, err error) {
	pwd, err := os.Getwd()

	if err != nil {
		log.Printf("Error: %s", err)
		return pdfRelativeLink, err
	}

	resultFolder := pwd + "/" + utils.Config.TempPath + "/result"

	if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
		if err := os.Mkdir(resultFolder, 0755); err != nil {
			log.Println("Cannot create result folder:", resultFolder, err)
			return pdfRelativeLink, err
		}
	}

	folderForResult := resultFolder + "/" + code
	relativeFolder := utils.Config.TempPath + "/result/" + code

	if _, err := os.Stat(folderForResult); os.IsNotExist(err) {
		if err := os.Mkdir(folderForResult, 0755); err != nil {
			log.Printf("Cannot create Folder: %s", err)
			return pdfRelativeLink, err
		}
	}

	pdfName := utils.Random()
	pdfPath := fmt.Sprintf("%s/%s.pdf", folderForResult, pdfName)
	pdfRelativeLink = fmt.Sprintf("%s/%s.pdf", relativeFolder, pdfName)

	err = utils.ImagesToPdf(pngsAbsolutePath, pdfPath)

	if err != nil {
		log.Printf("Error: %s\n", err)
		return pdfRelativeLink, err
	}

	return pdfRelativeLink, err
}
