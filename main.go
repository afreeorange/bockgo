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
		fmt.Print(version)
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
	v, _ := mem.VirtualMemory()

	articleRoot = strings.TrimRight(articleRoot, "/")
	outputFolder = strings.TrimRight(outputFolder, "/")
	statistics := Statistics{
		Architecture:   runtime.GOARCH,
		ArticleCount:   0,
		BuildDate:      time.Now().UTC(),
		CPUCount:       runtime.NumCPU(),
		GenerationTime: 0,
		MemoryInGB:     int(v.Total / (1024 * 1024 * 1024)),
		Platform:       runtime.GOOS,
	}

	config := BockConfig{
		articleRoot:  articleRoot,
		outputFolder: outputFolder,
		statistics:   statistics,
	}

	articles, err := process(config)
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

	config.statistics.GenerationTime = elapsed
	config.statistics.ArticleCount = len(articles)
	processHome(config)
}
