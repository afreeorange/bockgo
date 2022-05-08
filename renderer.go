package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var markdown = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
	goldmark.WithExtensions(
		extension.Footnote,
	),
)

//go:embed template
var templatesContent embed.FS

var pongoLoader = pongo2.NewFSLoader(templatesContent)
var templateSet = pongo2.NewSet("template", pongoLoader)

func renderArticle(source []byte, context Article) string {
	var conversionBuffer bytes.Buffer
	if err := markdown.Convert(source, &conversionBuffer); err != nil {
		panic(err)
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(context)
	json.Unmarshal(data, &myMap)

	t, _ := templateSet.FromCache("template/entity.html")
	o, _ := t.Execute(pongo2.Context{
		"id":          context.ID,
		"folder":      context.Folder,
		"hierarchy":   strings.Join(context.Path, "|"),
		"sizeInBytes": context.Size,
		"title":       context.Title,
		"type":        "article",
		"modified":    context.FileModified,
		"source":      context.Source,
		"html":        conversionBuffer.String(),
		"revisions":   context.Revisions,
		"version":     version,
	})

	conversionBuffer.Reset()

	return o
}

func renderHome(articleRoot string) {
	return
}
