package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

type ArticleFile struct {
	Folder       string   `json:"folder"`
	Path         []string `json:"path"`
	Size         int64    `json:"sizeInBytes"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Source       string   `json:"source"`
	Html         string   `json:"html"`
	FileModified string   `json:modified`
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func render(source []byte, buffer *bytes.Buffer) {
	if err := goldmark.Convert(source, buffer); err != nil {
		panic(err)
	}
}

func glob(articleRoot string, extension string) ([]ArticleFile, error) {
	files := []ArticleFile{}
	var buffer bytes.Buffer

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			dir := strings.Replace(filepath.Dir(path), articleRoot, "", -1)
			fileName := f.Name()
			title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			filePath := strings.Split(dir, "/")[1:]
			contents, _ := os.ReadFile(path)

			render(contents, &buffer)

			item := ArticleFile{
				Path:         filePath,
				Title:        title,
				Folder:       dir,
				Size:         f.Size(),
				Type:         "article",
				FileModified: f.ModTime().Format(time.RFC3339),
				Source:       string(contents),
				Html:         buffer.String(),
			}

			jsonData, _ := JSONMarshal(item)

			fmt.Println("Writing ", fileName)
			os.Mkdir("/Users/nikhilanand/Desktop/temp/"+title, os.ModePerm)
			os.WriteFile("/Users/nikhilanand/Desktop/temp/"+title+"/index.html", buffer.Bytes(), os.ModePerm)
			os.WriteFile("/Users/nikhilanand/Desktop/temp/"+title+"/index.json", jsonData, os.ModePerm)

			files = append(files, item)
			buffer.Reset()
		}

		return nil
	})

	return files, err
}

func main() {
	articleRoot := "/Users/nikhilanand/personal/wiki.nikhil.io.articles"

	list, err := glob(articleRoot, ".md")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("I found %s articles in %s\n", fmt.Sprint(len(list)), articleRoot)
}
