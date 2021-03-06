package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

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

func makeUri(path string, articleRoot string) string {
	uri := strings.ReplaceAll(strings.Replace(path, articleRoot, "", -1), " ", "_")
	return strings.TrimSuffix(uri, filepath.Ext(uri))
}

func makeID(path string) string {
	return uuid.NewV5(uuid.NamespaceURL, path).String()
}

func removeExtensionFrom(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func makeHierarchy(path string, articleRoot string) []Hierarchy {
	a := strings.Replace(path, articleRoot, "", -1)
	b := strings.Split(a, "/")
	c := []Hierarchy{}

	uriPath := ""

	for _, p := range b {
		uri := strings.ReplaceAll(strings.TrimSuffix(p, filepath.Ext(p)), " ", "_")

		if p == "" {
			c = append(c, Hierarchy{
				Name: "ROOT",
				Type: "folder",
				URI:  "/ROOT",
			})
		} else {
			name := strings.TrimSuffix(p, filepath.Ext(p))
			type_ := "folder"
			uriPath += "/" + uri
			uriPath = strings.TrimLeft(uriPath, "/")

			if filepath.Ext(p) == ".md" {
				type_ = "article"
			}

			c = append(c, Hierarchy{
				Name: name,
				Type: type_,
				URI:  uriPath,
			})
		}
	}

	return c
}
