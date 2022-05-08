package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// The JSON marshaller in Golang's STDLIB cannot be configured to disable HTML
// escaping. That's what this function does.
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func makeUri(articlePath string, articleRoot string) string {
	uri := strings.ReplaceAll(strings.Replace(articlePath, articleRoot, "", -1), " ", "_")
	return strings.TrimSuffix(uri, filepath.Ext(uri))
}

func makeID(articlePath string) string {
	return uuid.NewV5(uuid.NamespaceURL, articlePath).String()
}

func removeExtensionFrom(articlePath string) string {
	return strings.TrimSuffix(articlePath, filepath.Ext(articlePath))
}

func processArticle(
	articlePath string,
	articleRoot string,
	outputFolder string,
	f os.FileInfo,
) Article {
	dir := strings.Replace(filepath.Dir(articlePath), articleRoot, "", -1)
	fileName := f.Name()
	title := removeExtensionFrom(fileName)
	filePath := strings.Split(dir, "/")[1:]

	outPath := makeUri(articlePath, articleRoot)
	// fmt.Println("Writing", outPath)

	contents, _ := os.ReadFile(articlePath)
	item := Article{
		ID:           makeID(articlePath),
		Path:         filePath,
		Title:        title,
		Folder:       dir,
		Size:         f.Size(),
		FileModified: f.ModTime().UTC().Format(time.RFC3339),
		Source:       string(contents),
		Html:         "",
	}
	html := renderArticle(contents, item)
	item.Html = html

	os.MkdirAll(outputFolder+outPath, os.ModePerm)
	os.WriteFile(outputFolder+outPath+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(outputFolder+outPath+"/index.json", jsonData, os.ModePerm)

	return item
}

func process(
	articleRoot string,
	outputFolder string,
) ([]Article, error) {
	files := []Article{}

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if !IGNORED_FOLDERS_REGEX.MatchString(path) {

			if !IGNORED_FILES_REGEX.MatchString(path) {
				if filepath.Ext(path) == ".md" {
					item := processArticle(path, articleRoot, outputFolder, f)
					files = append(files, item)
				}
			}

			if f.IsDir() {
				if path == articleRoot {
					fmt.Println("Will render home")
				} else {
					fmt.Println("Found", path)
				}
			}
		}

		return nil
	})

	return files, err
}
