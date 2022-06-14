package main

import (
	"embed"
	"regexp"
	"strings"

	chroma "github.com/alecthomas/chroma/formatters/html"
	"github.com/flosch/pongo2/v5"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var VERSION = "0.0.6-alpha"
var DATE_LAYOUT = "2006-01-02 15:04:05 -0700"
var DATABASE_NAME = "articles.db"

// Exit codes
var EXIT_NO_ARTICLE_ROOT = 1
var EXIT_NO_OUTPUT_FOLDER = 2
var EXIT_NOT_A_GIT_REPO = 3

// TODO: Ignore dotfolders
var IGNORED_FOLDERS_REGEX = regexp.MustCompile(
	strings.Join([]string{
		"__assets",
		"_assets",
		"\\.circleci",
		"\\.git",
		"css",
		"img",
		"js",
		"node_modules",
	}, "|"))
var IGNORED_FILES_REGEX = regexp.MustCompile("Home.md")

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
