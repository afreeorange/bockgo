package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func processArticle(
	articlePath string,
	articleRoot string,
	outputFolder string,
	f os.FileInfo,
) Article {
	// dir := strings.Replace(filepath.Dir(articlePath), articleRoot, "", -1)
	// filePath := strings.Split(dir, "/")[1:]

	fileName := f.Name()
	title := removeExtensionFrom(fileName)
	uri := makeUri(articlePath, articleRoot)

	// fmt.Println("Writing", uri)

	contents, _ := os.ReadFile(articlePath)
	item := Article{
		ID:           makeID(articlePath),
		URI:          uri,
		Title:        title,
		Size:         f.Size(),
		FileModified: f.ModTime().UTC().Format(time.RFC3339),
		Source:       string(contents),
		Html:         "",
	}
	html := renderArticle(contents, item, "article")
	item.Html = html

	os.MkdirAll(outputFolder+uri, os.ModePerm)
	os.WriteFile(outputFolder+uri+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(outputFolder+uri+"/index.json", jsonData, os.ModePerm)

	return item
}

func processHome(
	articleRoot string,
	outputFolder string,
) {
	homePath := articleRoot + "/Home.md"
	var contents []byte
	var size int64
	var mTime string

	f, err := os.Stat(homePath)

	if err != nil {
		fmt.Println("Could not find Home.md... making one.")

		contents = []byte("(You need to make a `Home.md` here!)")
		mTime = time.Now().UTC().Format(time.RFC3339)
		size = 0
	} else {
		g, _ := os.ReadFile(articleRoot + "/Home.md")

		contents = g
		mTime = f.ModTime().UTC().Format(time.RFC3339)
		size = f.Size()
	}

	item := Article{
		ID:           makeID(articleRoot + "/Home.md"),
		Title:        "Hello!",
		Size:         size,
		FileModified: mTime,
		Source:       string(contents),
		Html:         "",
		URI:          "/",
	}
	html := renderArticle(contents, item, "home")
	item.Html = html

	os.WriteFile(outputFolder+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := jsonMarshal(item)
	os.WriteFile(outputFolder+"/index.json", jsonData, os.ModePerm)
}

func processArchive(articles []Article, outputFolder string) {
	os.MkdirAll(outputFolder+"/archive", os.ModePerm)
	html := renderArchive(articles)
	os.WriteFile(outputFolder+"/archive/index.html", []byte(html), os.ModePerm)
}

func process(
	articleRoot string,
	outputFolder string,
) ([]Article, error) {
	files := []Article{}

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if !IGNORED_FOLDERS_REGEX.MatchString(path) {

			if !IGNORED_FILES_REGEX.MatchString(path) {
				if filepath.Ext(path) == ".md" {
					item := processArticle(path, articleRoot, outputFolder, f)
					files = append(files, item)
				}
			}

			// This is a folder
			if f.IsDir() {
				if path != articleRoot {
					fmt.Println("Found", path)
				}
			}

			// The Homepage is special and will be processed separately.
		}

		return nil
	})

	return files, err
}
