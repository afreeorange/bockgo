package main

import (
	"bytes"
	"encoding/json"

	"github.com/flosch/pongo2/v5"
)

func renderArticle(source []byte, context Article, entityType string, config BockConfig) string {
	var conversionBuffer bytes.Buffer
	if err := markdown.Convert(source, &conversionBuffer); err != nil {
		panic(err)
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(context)
	json.Unmarshal(data, &myMap)

	t, _ := templateSet.FromCache("template/article.html")
	o, _ := t.Execute(pongo2.Context{
		"id":          context.ID,
		"sizeInBytes": context.Size,
		"title":       context.Title,
		"modified":    context.FileModified,
		"source":      context.Source,
		"uri":         context.URI,
		"html":        conversionBuffer.String(),
		"hierarchy":   context.Hierarchy,

		"type":       entityType,
		"version":    version,
		"statistics": config.statistics,
	})

	conversionBuffer.Reset()

	return o
}

func renderArchive(articles []Article) string {
	t, _ := templateSet.FromCache("template/archive.html")
	o, _ := t.Execute(pongo2.Context{
		"title":    "Archive",
		"articles": articles,

		"type":    "archive",
		"version": version,
	})

	return o
}
