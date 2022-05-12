package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func processArticle(articlePath string, config BockConfig, f os.FileInfo) Article {
	fileName := f.Name()
	title := removeExtensionFrom(fileName)
	uri := makeUri(articlePath, config.articleRoot)

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
	html, raw := renderArticle(contents, item, "article", config)
	item.Html = html

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
	os.WriteFile(config.outputFolder+"/Home/index.html", []byte(html), os.ModePerm)

	os.MkdirAll(config.outputFolder+"/Home/raw", os.ModePerm)
	os.WriteFile(config.outputFolder+"/raw/index.html", []byte(raw), os.ModePerm)
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
		if f.IsDir() {
			folders = append(folders, FolderThing{
				Name: f.Name(),
				Type: "folder",
				Path: "a",
				URI:  makeUri(path, config.articleRoot),
			})
		} else {
			articles = append(articles, FolderThing{
				Name: f.Name(),
				Type: "article",
				Path: "a",
				URI:  makeUri(path, config.articleRoot),
			})
		}
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
		README:    "",
	}

	// d, _ := json.MarshalIndent(context, "", "  ")
	// fmt.Print(string(d))

	html := renderFolder(context)

	os.MkdirAll(config.outputFolder+"/"+makeUri(path, config.articleRoot), os.ModePerm)
	os.WriteFile(config.outputFolder+"/"+makeUri(path, config.articleRoot)+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := jsonMarshal(context)
	os.WriteFile(config.outputFolder+"/"+makeUri(path, config.articleRoot)+"/index.json", jsonData, os.ModePerm)

	return folders, articles
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
					processFolder(path, config)
				}
			}

			// The Homepage is special and will be processed separately.
		}

		return nil
	})

	return files, err
}
