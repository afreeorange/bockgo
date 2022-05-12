package main

import (
	"bytes"
	"encoding/json"

	"github.com/flosch/pongo2/v5"
)

func renderArticle(
	source []byte,
	context Article,
	entityType string,
	config BockConfig,
) (string, string) {
	var conversionBuffer bytes.Buffer
	if err := markdown.Convert(source, &conversionBuffer); err != nil {
		panic(err)
	}

	var myMap map[string]interface{}
	data, _ := json.Marshal(context)
	json.Unmarshal(data, &myMap)

	baseContext := pongo2.Context{
		"id":          context.ID,
		"sizeInBytes": context.Size,
		"title":       context.Title,
		"modified":    context.FileModified,
		"source":      context.Source,
		"uri":         context.URI,
		"html":        conversionBuffer.String(),
		"hierarchy":   context.Hierarchy,

		"version":    version,
		"statistics": config.statistics,
	}

	baseContext.Update(pongo2.Context{
		"type": "article",
	})

	t1, _ := templateSet.FromCache("template/article.html")
	o1, _ := t1.Execute(baseContext)

	baseContext.Update(pongo2.Context{
		"type": "raw",
	})

	t2, _ := templateSet.FromCache("template/raw.html")
	o2, _ := t2.Execute(baseContext)

	conversionBuffer.Reset()

	return o1, o2
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

func renderFolder(context Folder) string {
	t, _ := templateSet.FromCache("template/folder.html")
	o, _ := t.Execute(pongo2.Context{
		"title":     context.Title,
		"children":  context.Children,
		"hierarchy": context.Hierarchy,
		"readme":    context.README,

		"type":    "folder",
		"version": version,
	})

	return o
}
