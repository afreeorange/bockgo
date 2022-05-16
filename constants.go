package main

import (
	"embed"
	"regexp"

	chroma "github.com/alecthomas/chroma/formatters/html"
	"github.com/flosch/pongo2/v5"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// TODO: Ignore dotfolders
var IGNORED_FOLDERS_REGEX = regexp.MustCompile("node_modules|.git|.circleci|_assets|js|css|img")

var IGNORED_FILES_REGEX = regexp.MustCompile("Home.md")

var DATABASE_NAME = "articles.db"

var markdown = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	),
	goldmark.WithExtensions(
		extension.Footnote,
		extension.Linkify,
		extension.Strikethrough,
		extension.Table,
		extension.Typographer,
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chroma.WithClasses(true),
			),
		),
	),
)

//go:embed template
var templatesContent embed.FS

var pongoLoader = pongo2.NewFSLoader(templatesContent)
var templateSet = pongo2.NewSet("template", pongoLoader)
