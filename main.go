package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

//go:embed VERSION
var b []byte
var version = string(b)

func main() {
	var versionInfo bool
	var articleRoot string
	var outputFolder string

	flag.BoolVar(&versionInfo, "v", false, "Version info")
	flag.StringVar(&articleRoot, "a", "", "Article root")
	flag.StringVar(&outputFolder, "o", "", "Output folder")

	flag.Parse()

	if versionInfo {
		fmt.Println(version)
		os.Exit(0)
	}

	if articleRoot == "" {
		fmt.Println("You must give me an article root")
		os.Exit(1)
	}

	if outputFolder == "" {
		fmt.Println("You must give me an output folder")
		os.Exit(2)
	}

	if _, err := git.PlainOpen(articleRoot); err != nil {
		fmt.Println("That article root does not appear to be a git repository.")
		os.Exit(3)
	}

	articleRoot = strings.TrimRight(articleRoot, "/")
	outputFolder = strings.TrimRight(outputFolder, "/")

	_, err := process(articleRoot, outputFolder, ".md")
	if err != nil {
		fmt.Println("Error occurred processing repository:", err)
	}

	copyAssets(articleRoot, outputFolder)
}
