package main

import "regexp"

var IGNORED_FOLDERS_REGEX = regexp.MustCompile("node_modules|.git|.circleci|_assets|js|css|img")

var IGNORED_FILES_REGEX = regexp.MustCompile("Home.md|README.md")

var TEMPLATE_ASSETS = [...]string{
	"styles.css",
	"robots.txt",
}
