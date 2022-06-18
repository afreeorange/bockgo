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
	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	var versionInfo bool
	var articleRoot string
	var outputFolder string
	var makeJSON bool

	flag.BoolVar(&versionInfo, "v", false, "Version info")
	flag.StringVar(&articleRoot, "a", "", "Article root")
	flag.StringVar(&outputFolder, "o", "", "Output folder")
	flag.BoolVar(&makeJSON, "j", true, "Create JSON source files")

	flag.Parse()

	if versionInfo {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if articleRoot == "" {
		fmt.Println("You must give me an article root.")
		os.Exit(EXIT_NO_ARTICLE_ROOT)
	}

	if outputFolder == "" {
		fmt.Println("You must give me an output folder.")
		os.Exit(EXIT_NO_OUTPUT_FOLDER)
	}

	// Some bookkeeping
	start := time.Now()
	v, _ := mem.VirtualMemory()

	repository, err := git.PlainOpen(articleRoot)

	if err != nil {
		fmt.Println("That article root does not appear to be a git repository.")
		os.Exit(EXIT_NOT_A_GIT_REPO)
	}

	articleRoot = strings.TrimRight(articleRoot, "/")
	outputFolder = strings.TrimRight(outputFolder, "/")

	// App config
	config := BockConfig{
		articleRoot:  articleRoot,
		outputFolder: outputFolder,
		meta: Meta{
			Architecture:   runtime.GOARCH,
			ArticleCount:   0,
			BuildDate:      time.Now().UTC(),
			CPUCount:       runtime.NumCPU(),
			GenerationTime: 0,
			MemoryInGB:     int(v.Total / (1024 * 1024 * 1024)),
			Platform:       runtime.GOOS,
		},
		makeJSON: makeJSON,
		started:  time.Now(),
	}

	// Create the output folder first
	fmt.Println("Making", outputFolder, "if it doesn't exist")
	os.MkdirAll(outputFolder, os.ModePerm)

	// Database setup
	db := makeDatabase(config)
	defer db.Close()

	articles, err := process(config, repository, db)
	if err != nil {
		fmt.Println("Could not process article root: ", err)
	}

	processArchive(articles, config)
	copyAssets(config)
	copyTemplateAssets(config)

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf(
		"Finished processing %d articles in %s\n",
		len(articles),
		time.Duration.Round(elapsed, time.Millisecond),
	)

	// Generate the home page once all the statistics are gathered
	config.meta.GenerationTime = elapsed
	config.meta.ArticleCount = len(articles)
	processHome(config)
}
