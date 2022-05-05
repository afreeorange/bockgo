package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func json_marshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func process(
	articleRoot string,
	outputFolder string,
	extension string,
	r *git.Repository,
) ([]Article, error) {
	files := []Article{}

	messages := make(chan Article)

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			go processArticle(path, articleRoot, outputFolder, f, r, messages)

			item := <-messages
			files = append(files, item)
		}

		return nil
	})

	return files, err
}
