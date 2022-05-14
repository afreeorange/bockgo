package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

func getCommits(repository *git.Repository, fileName string, config BockConfig, revisionsChannel chan []Revision) []Revision {
	o, err := exec.Command(
		"git",
		"-C",
		config.articleRoot,
		"log",
		`--pretty=format:'{"id": "%H", "shortId": "%h", "subject": "%f", "body": "%b", "date": "%aD"}'`,
		fileName,
	).Output()

	res := []Revision{}

	if err != nil {
		fmt.Println("Error getting commits:", err)
		return res
	} else {
		s := string(o)
		s = strings.ReplaceAll(s, "'\n'", ",")
		s = strings.ReplaceAll(s, "'", "")
		s = strings.ReplaceAll(s, "\n", "")
		s = "[" + s + "]"

		json.Unmarshal([]byte(s), &res)
	}

	revisionsChannel <- res

	return res
}

func processArticle(articlePath string, config BockConfig, f os.FileInfo, repository *git.Repository, stmt *sql.Stmt) Article {
	fileName := f.Name()
	title := removeExtensionFrom(fileName)
	uri := makeUri(articlePath, config.articleRoot)
	// revisionsChannel := make(chan []Revision)

	contents, _ := os.ReadFile(articlePath)
	item := Article{
		ID:           makeID(articlePath),
		URI:          uri,
		Title:        title,
		Size:         f.Size(),
		FileModified: f.ModTime().UTC(),
		Source:       string(contents),
		Html:         "",
		Hierarchy:    makeHierarchy(articlePath, config.articleRoot),
		Revisions:    nil,
	}

	// Insert into Database
	_, se := stmt.Exec(makeID(articlePath), string(contents), f.ModTime().UTC(), title, uri)
	if se != nil {
		log.Fatal("SHIT ", se)
	}

	// Render the article HTML
	html, raw := renderArticle(contents, item, "article", config)
	item.Html = html

	// fmt.Println(f.Name())
	// go getCommits(repository, articlePath, config, revisionsChannel)
	// revisions := <-revisionsChannel
	// item.Revisions = revisions

	os.MkdirAll(config.outputFolder+uri+"/raw", os.ModePerm)
	os.WriteFile(config.outputFolder+uri+"/index.html", []byte(html), os.ModePerm)
	os.WriteFile(config.outputFolder+uri+"/raw/index.html", []byte(raw), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(config.outputFolder+uri+"/index.json", jsonData, os.ModePerm)

	return item
}

func processHome(config BockConfig) {
	homePath := config.articleRoot + "/Home.md"
	var contents []byte
	var size int64
	var mTime time.Time

	f, err := os.Stat(homePath)

	if err != nil {
		fmt.Println("Could not find Home.md... making one.")

		contents = []byte("(You need to make a `Home.md` here!)")
		mTime = time.Now().UTC()
		size = 0
	} else {
		g, _ := os.ReadFile(config.articleRoot + "/Home.md")

		contents = g
		mTime = f.ModTime().UTC()
		size = f.Size()
	}

	item := Article{
		ID:           makeID(config.articleRoot + "/Home.md"),
		Title:        "Hello!",
		Size:         size,
		FileModified: mTime,
		Source:       string(contents),
		Html:         "",
		URI:          "",
		Hierarchy:    makeHierarchy("/Home.md", config.articleRoot),
	}
	html, raw := renderArticle(contents, item, "home", config)
	item.Html = html

	os.MkdirAll(config.outputFolder+"/raw", os.ModePerm)
	os.WriteFile(config.outputFolder+"/index.html", []byte(html), os.ModePerm)
	os.WriteFile(config.outputFolder+"/raw/index.html", []byte(raw), os.ModePerm)

	os.MkdirAll(config.outputFolder+"/Home/raw", os.ModePerm)
	os.WriteFile(config.outputFolder+"/Home/index.html", []byte(html), os.ModePerm)
	os.WriteFile(config.outputFolder+"/Home/raw/index.html", []byte(raw), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(config.outputFolder+"/index.json", jsonData, os.ModePerm)
}

func processArchive(articles []Article, config BockConfig) {
	os.MkdirAll(config.outputFolder+"/archive", os.ModePerm)
	html := renderArchive(articles)
	os.WriteFile(config.outputFolder+"/archive/index.html", []byte(html), os.ModePerm)
}

func processFolder(path string, config BockConfig) ([]FolderThing, []FolderThing) {
	l, _ := ioutil.ReadDir(path)

	folders := []FolderThing{}
	articles := []FolderThing{}
	title := strings.TrimLeft(strings.Replace(path, config.articleRoot, "", -1), "/")

	for _, f := range l {
		if !IGNORED_FOLDERS_REGEX.MatchString(f.Name()) {
			if f.IsDir() {
				folders = append(folders, FolderThing{
					Name: removeExtensionFrom(f.Name()),
					Type: "folder",
					URI:  makeUri(path, config.articleRoot) + "/" + makeUri(f.Name(), config.articleRoot),
				})
			} else {
				articles = append(articles, FolderThing{
					Name: removeExtensionFrom(f.Name()),
					Type: "article",
					URI:  makeUri(path, config.articleRoot) + "/" + makeUri(f.Name(), config.articleRoot),
				})
			}
		}
	}

	// Check if the folder has a readme
	readme := ""
	r, err := os.ReadFile(path + "/README.md")
	if err == nil {
		readme = string(r)
	}

	context := Folder{
		ID:    makeID(path),
		URI:   makeUri(path, config.articleRoot),
		Title: title,
		Children: Children{
			Articles: articles,
			Folders:  folders,
		},
		Hierarchy: makeHierarchy(path, config.articleRoot),
		README:    readme,
	}

	html := renderFolder(context)

	if path != config.articleRoot {
		os.MkdirAll(config.outputFolder+"/"+makeUri(path, config.articleRoot), os.ModePerm)
		os.WriteFile(config.outputFolder+"/"+makeUri(path, config.articleRoot)+"/index.html", []byte(html), os.ModePerm)

		jsonData, _ := jsonMarshal(context)
		os.WriteFile(config.outputFolder+"/"+makeUri(path, config.articleRoot)+"/index.json", jsonData, os.ModePerm)
	} else {
		os.MkdirAll(config.outputFolder+"/ROOT", os.ModePerm)
		os.WriteFile(config.outputFolder+"/ROOT/index.html", []byte(html), os.ModePerm)

		jsonData, _ := jsonMarshal(context)
		os.WriteFile(config.outputFolder+"/ROOT/index.json", jsonData, os.ModePerm)
	}

	return folders, articles
}

func process(config BockConfig, repository *git.Repository, db *sql.DB) ([]Article, error) {
	files := []Article{}
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`
  INSERT INTO articles (
      id,
      content,
      modified,
      title,
      uri
    )
    VALUES (?, ?, ?, ?, ?)
  `)

	defer stmt.Close()
	err := filepath.Walk(config.articleRoot, func(path string, f os.FileInfo, err error) error {
		if !IGNORED_FOLDERS_REGEX.MatchString(path) {

			if !IGNORED_FILES_REGEX.MatchString(path) {
				if filepath.Ext(path) == ".md" {
					item := processArticle(path, config, f, repository, stmt)
					files = append(files, item)
				}
			}

			// This is a folder
			if f.IsDir() {
				processFolder(path, config)
			}

			// The Homepage is special and will be processed separately.
		}

		return nil
	})

	tx.Commit()

	return files, err
}
