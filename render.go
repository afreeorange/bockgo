package main

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func render(source []byte, buffer *bytes.Buffer) {
	if err := goldmark.Convert(source, buffer); err != nil {
		panic(err)
	}
}
