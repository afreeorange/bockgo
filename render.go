package main

import (
	"bytes"
	"embed"
	"html/template"

	"github.com/yuin/goldmark"
)

//go:embed *.tmpl
var content embed.FS
var t, _ = template.ParseFS(content, []string{
	"base.tmpl",
	"entity.tmpl",
}...)

func render(source []byte, context Article) string {
	var conversionBuffer bytes.Buffer
	if err := goldmark.Convert(source, &conversionBuffer); err != nil {
		panic(err)
	}

	var outputBuffer bytes.Buffer
	t.ExecuteTemplate(
		&outputBuffer,
		"base",
		struct {
			Context Article
			HTML    template.HTML
		}{
			Context: context,
			HTML:    template.HTML(conversionBuffer.String()),
		},
	)
	o := outputBuffer.String()

	conversionBuffer.Reset()
	outputBuffer.Reset()

	return o
}
