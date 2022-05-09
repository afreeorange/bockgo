package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/shirou/gopsutil/v3/mem"
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
		fmt.Println("You must give me an article root.")
		os.Exit(1)
	}

	if outputFolder == "" {
		fmt.Println("You must give me an output folder.")
		os.Exit(2)
	}

	if _, err := git.PlainOpen(articleRoot); err != nil {
		fmt.Println("That article root does not appear to be a git repository.")
		os.Exit(3)
	}

	// Some bookkeeping
	start := time.Now()

	articleRoot = strings.TrimRight(articleRoot, "/")
	outputFolder = strings.TrimRight(outputFolder, "/")

	articles, err := process(articleRoot, outputFolder)
	if err != nil {
		fmt.Println("Could not process article root: ", err)
	}

	// https://stackoverflow.com/questions/28999735/what-is-the-shortest-way-to-simply-sort-an-array-of-structs-by-arbitrary-field
	// sortedArticles, _ := sort.Slice(articles, func(i, j int) bool {
	// 	return articles[i].FileModified > articles[j].FileModified
	// })
	// for _, a := range articles {
	// 	fmt.Println(a.Title, a.FileModified)
	// }

	processArchive(articles, outputFolder)
	copyAssets(articleRoot, outputFolder)
	copyTemplateAssets(outputFolder)

	end := time.Now()
	v, _ := mem.VirtualMemory()

	statistics := Statistics{
		GenerationTime: end.Sub(start),
		ArticleCount:   len(articles),
		CPUs:           runtime.NumCPU(),
		Memory:         int(v.Total),
		Platform:       runtime.GOOS,
		Architecture:   runtime.GOARCH,
	}
	_ = statistics

	fmt.Printf("Finished processing %d articles in %s", len(articles), end.Sub(start))
	processHome(articleRoot, outputFolder)
}
