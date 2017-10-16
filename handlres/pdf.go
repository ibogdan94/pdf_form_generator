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
	"bytes"
	"image/png"
	"log"
	"encoding/base64"
	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

type B64 struct {
	B64 string `form:"b64" json:"b64" binding:"required"`
}

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

	if len(result) == 0 {
		c.Redirect(301, "/")
	}

	fmt.Println("Inside len")

	session, err := store.Get(c.Request, "generated_imgs")

	if err != nil {
		fmt.Println("Session error:", err)
		c.Redirect(301, "/")
	}

	session.Values["images"] = result

	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Println(err)
		c.Redirect(301, "/")
	}

	c.Redirect(301, "/pdf/edit")
}

func EditParsedPDF(c *gin.Context) {
	gfSession, err := store.Get(c.Request, "generated_imgs")

	//no one images
	if err != nil {
		c.Redirect(301, "/")
	}

	imgs := gfSession.Values
	fmt.Println(imgs)


	c.HTML(http.StatusOK, "builder.html", imgs)
}

func SavePDF(ctx *gin.Context) {
	var json B64

	if err := ctx.BindJSON(&json); err != nil {
		fmt.Println("Json error:", err)
	}

	var b64FromRequest string = json.B64
	var b64 string = b64FromRequest[22:len(b64FromRequest)]

	unbased, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		fmt.Println("Cannot decode b64. Error:", err)
		panic(err)
	}

	r := bytes.NewReader(unbased)
	img, err := png.Decode(r)

	if err != nil {
		panic("Bad png")
	}

	folderForResult := "./temp/result"

	if stat, err := os.Stat(folderForResult); err != nil && stat.IsDir() == false {
		if err := os.Mkdir(folderForResult, 0755); err != nil {
			fmt.Println("Cannot Create Folder:", err)
			log.Fatal(err)
		}
	}


	fileName := utils.Random()
	targetPath := folderForResult + "/" + fileName + ".png"

	f, err := os.Create(targetPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	pdfPath := folderForResult + "/" + fileName + ".pdf"

	//pdf := gofpdf.New("P", "mm", "A4", "")
	//pdf.ImageOptions(targetPath, 0, 0, 210, 297, true, gofpdf.ImageOptions{"PNG", true}, 0, "")

	//if err := pdf.OutputFileAndClose(pdfPath); err != nil {
	//	fmt.Println("Cannot create PDF by image:", err)
	//}

	inputPaths := []string{targetPath}

	err = imagesToPdf(inputPaths, pdfPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"url": pdfPath[1:],
	})
}

func imagesToPdf(inputPaths []string, outputPath string) error {
	c := creator.New()

	for _, imgPath := range inputPaths {
		unicommon.Log.Debug("Image: %s", imgPath)

		img, err := creator.NewImageFromFile(imgPath)
		if err != nil {
			unicommon.Log.Debug("Error loading image: %v", err)
			return err
		}
		img.ScaleToWidth(612.0)

		height := 612.0 * img.Height() / img.Width()
		c.SetPageSize(creator.PageSize{612, height})
		c.NewPage()
		img.SetPos(0, 0)
		err = c.Draw(img)

		if err != nil {
			return err
		}
	}

	err := c.WriteToFile(outputPath)
	return err
}

func parsePDFToPng(file multipart.File, headers *multipart.FileHeader) (result []string, error error) {
	randomFileName := utils.Random()
	folderForPDF := "./temp/" + randomFileName

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

	//fmt.Println("numberOfPages", numberOfPages)

	if numberOfPages > 0 {
		//fmt.Println("Create wait group")
		wg := new(sync.WaitGroup)

		for i := 0; i < int(numberOfPages); i++ {
			page := fmt.Sprintf("%s[%d]", imageFilePathWithExtension, i)
			output := fmt.Sprintf("%s/%s.%d.png", folderForPDF, utils.Random(), i)
			//fmt.Println("------------------")
			//fmt.Println(output)
			//fmt.Println("------------------")
			command := fmt.Sprintf("convert -verbose -trim -density %d -resize %s -depth %d -flatten %s %s", 400, "25%", 8, page, output)
			fmt.Println(command)
			wg.Add(1)
			result = append(result, output[1:])
			go exe_cmd(command, wg)
		}

		wg.Wait()
	}

	//fmt.Println("PNG RESULT")
	//fmt.Println(result)

	return result, err
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
