package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	_ "github.com/mattn/go-sqlite3"
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

	repository, err := git.PlainOpen(articleRoot)

	if err != nil {
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

	// Set up the database
	dbPath := config.outputFolder + "/articles.db"
	os.Remove(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS articles (
    id              TEXT NOT NULL UNIQUE,
    content         TEXT,
    modified        TEXT NOT NULL,
    title           TEXT NOT NULL,
    uri             TEXT NOT NULL
  );
  CREATE VIRTUAL TABLE articles_fts USING fts5(
    id,
    content,
    modified,
    title,
    uri,
    content="articles"
  );
  CREATE TRIGGER fts_update AFTER INSERT ON articles
    BEGIN
      INSERT INTO articles_fts (
        id,
        content,
        modified,
        title,
        uri
      )
      VALUES (
        new.id,
        new.content,
        new.modified,
        new.title,
        new.uri
      );
  END;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// --- End Database Setup ---

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

	config.statistics.GenerationTime = elapsed
	config.statistics.ArticleCount = len(articles)
	processHome(config)
}
