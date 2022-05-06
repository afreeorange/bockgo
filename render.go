package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/yuin/goldmark"
)

//go:embed templates
var content embed.FS

var pongoLoader = pongo2.NewFSLoader(content)
var templateSet = pongo2.NewSet("templates", pongoLoader)

func render(source []byte, context Article) string {
	var conversionBuffer bytes.Buffer
	if err := goldmark.Convert(source, &conversionBuffer); err != nil {
		panic(err)
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(context)
	json.Unmarshal(data, &myMap)

	t, _ := templateSet.FromCache("templates/entity.html")
	o, _ := t.Execute(pongo2.Context{
		"id":          context.ID,
		"folder":      context.Folder,
		"path":        strings.Join(context.Path, "|"),
		"sizeInBytes": context.Size,
		"title":       context.Title,
		"type":        context.Type,
		"modified":    context.FileModified,
		"source":      context.Source,
		"html":        conversionBuffer.String(),
		"revisions":   context.Revisions,
	})

	conversionBuffer.Reset()

	return o
}
