package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func processArticle(articlePath string, config BockConfig, f os.FileInfo) Article {
	// dir := strings.Replace(filepath.Dir(articlePath), articleRoot, "", -1)
	// filePath := strings.Split(dir, "/")[1:]

	fileName := f.Name()
	title := removeExtensionFrom(fileName)
	uri := makeUri(articlePath, config.articleRoot)

	// fmt.Println("Writing", uri)

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
	}
	html := renderArticle(contents, item, "article", config)
	item.Html = html

	os.MkdirAll(config.outputFolder+uri, os.ModePerm)
	os.WriteFile(config.outputFolder+uri+"/index.html", []byte(html), os.ModePerm)

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
	html := renderArticle(contents, item, "home", config)
	item.Html = html

	os.WriteFile(config.outputFolder+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(config.outputFolder+"/index.json", jsonData, os.ModePerm)
}

func processArchive(articles []Article, config BockConfig) {
	os.MkdirAll(config.outputFolder+"/archive", os.ModePerm)
	html := renderArchive(articles)
	os.WriteFile(config.outputFolder+"/archive/index.html", []byte(html), os.ModePerm)
}

func process(config BockConfig) ([]Article, error) {
	files := []Article{}

	err := filepath.Walk(config.articleRoot, func(path string, f os.FileInfo, err error) error {
		if !IGNORED_FOLDERS_REGEX.MatchString(path) {

			if !IGNORED_FILES_REGEX.MatchString(path) {
				if filepath.Ext(path) == ".md" {
					item := processArticle(path, config, f)
					files = append(files, item)
				}
			}

			// This is a folder
			if f.IsDir() {
				if path != config.articleRoot {
					fmt.Println("Found", path)
				}
			}

			// The Homepage is special and will be processed separately.
		}

		return nil
	})

	return files, err
}
