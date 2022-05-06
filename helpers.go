package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
	uuid "github.com/satori/go.uuid"
)

func makeUri(articlePath string, articleRoot string) string {
	uri := strings.ReplaceAll(strings.Replace(articlePath, articleRoot, "", -1), " ", "_")
	return strings.TrimSuffix(uri, filepath.Ext(uri))
}

func copyAssets(articleRoot string, outputFolder string) {
	fmt.Print("Copying assets... ")

	err := cp.Copy(articleRoot+"/__assets", outputFolder+"/assets")
	if err != nil {
		fmt.Println("Could not copy assets: ", err)
	}

	fmt.Println("done.")
}

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

func processArticle(
	articlePath string,
	articleRoot string,
	outputFolder string,
	f os.FileInfo,
) Article {
	dir := strings.Replace(filepath.Dir(articlePath), articleRoot, "", -1)
	fileName := f.Name()
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	filePath := strings.Split(dir, "/")[1:]

	outPath := makeUri(articlePath, articleRoot)
	fmt.Println("Writing", outPath)

	contents, _ := os.ReadFile(articlePath)
	item := Article{
		ID:           uuid.NewV5(uuid.NamespaceURL, articlePath).String(),
		Path:         filePath,
		Title:        title,
		Folder:       dir,
		Size:         f.Size(),
		Type:         "article",
		FileModified: f.ModTime().UTC().Format(time.RFC3339),
		Source:       string(contents),
		Html:         "",
	}
	html := render(contents, item)
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
	extension string,
) ([]Article, error) {
	files := []Article{}

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			item := processArticle(path, articleRoot, outputFolder, f)
			files = append(files, item)
		}

		return nil
	})

	return files, err
}
