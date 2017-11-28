package utils

import (
	"os"
	"path/filepath"
	"log"
	"fmt"
)

func Walker(tempPath string, codes []string) (err error)  {
	resultFolder := tempPath + "/result"

	if stat, err := os.Stat(resultFolder); err == nil && stat.IsDir() {
		//the fastest way is remove result folder and create new one
		if err := removeFolderWithFileInside(resultFolder); err != nil {
			log.Println(err)
			return nil
		}

		if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
			if err := os.Mkdir(resultFolder, 0755); err != nil {
				log.Println("Cannot create result folder:", resultFolder, err)
				return err
			}
		}
	}

	err = filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			save := false
			for _, v := range codes {
				if v == info.Name() || info.Name() == "temp" || info.Name() == "result" {
					save = true
				}
			}

			if save == false {
				folder := tempPath + "/" + info.Name()
				fmt.Println(info.Name())
				if err := removeFolderWithFileInside(folder); err != nil {
					log.Println(err)
					return err
				}
			}

			return nil
		}
		return nil
	})

	return err
}

func removeFolderWithFileInside(folder string) (err error)  {
	if stat, err := os.Stat(folder); err == nil && stat.IsDir() {
		if err := os.RemoveAll(folder); err != nil {
			return err
		}
	}
	return err
}
