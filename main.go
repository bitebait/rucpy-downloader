package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/gocolly/colly"
)

var destPath *string

func init() {
	destPath = flag.String("d", "", "Directory where the files will be saved.")
}

func main() {
	flag.Parse()
	path := *destPath
	if path == "" {
		log.Fatal("Please specify the directory where the files will be saved.")
	} else {
		getFiles(*destPath)
		log.Println("Finished.\nFiles saved in: ", *destPath)
	}
}

func unzipFile(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func crawler() []string {
	var filesUrl []string

	url := "https://www.set.gov.py/portal/PARAGUAY-SET/InformesPeriodicos?folder-id=repository:collaboration:/sites/PARAGUAY-SET/categories/SET/Informes%20Periodicos/listado-de-ruc-con-sus-equivalencias"

	c := colly.NewCollector(colly.Async(true))

	c.OnHTML("div.uiContentBox div.heading a", func(e *colly.HTMLElement) {
		fileName := e.Text
		if fileName != "" {
			filesUrl = append(filesUrl, "https://www.set.gov.py/rest/contents/download/collaboration/sites/PARAGUAY-SET/documents/informes-periodicos/ruc/"+fileName)
			log.Println("File found:", fileName)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting crawler...")
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit(url)

	c.Wait()

	return filesUrl
}

func downloadFiles(destPath string, filesUrl []string) []string {
	var files []string

	for _, url := range filesUrl {
		file, err := grab.Get(destPath, url)
		if err != nil {
			log.Fatalf("Error downloading the file: %s\nError: %s", url, err)
		}

		files = append(files, file.Filename)
		log.Println("File successfully downloaded:", file.Filename)
	}

	return files
}

func getFiles(destDir string) {
	checkDir(destDir) // check if dir exists.

	fileNames := crawler()

	files := downloadFiles(destDir, fileNames)

	for _, file := range files {
		log.Println("Extracting file:", file)
		err := unzipFile(file, destDir)
		if err != nil {
			log.Fatal("Error extracting the file:", file)
		}
		os.Remove(file)
	}
}

func checkDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Directory not found: %s\nPlease enter a valid directory.", dir)
	}
}
