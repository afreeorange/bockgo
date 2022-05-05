package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	cp "github.com/otiai10/copy"
	uuid "github.com/satori/go.uuid"
)

func makeUri(articlePath string, articleRoot string) string {
	uri := strings.ReplaceAll(strings.Replace(articlePath, articleRoot, "", -1), " ", "_")
	return strings.TrimSuffix(uri, filepath.Ext(uri))
}

func processArticle(
	articlePath string,
	articleRoot string,
	outputFolder string,
	f os.FileInfo,
	repository *git.Repository,
	messages chan Article,
) {
	dir := strings.Replace(filepath.Dir(articlePath), articleRoot, "", -1)
	fileName := f.Name()
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	filePath := strings.Split(dir, "/")[1:]

	fmt.Println("Processing", makeUri(articlePath, articleRoot))

	// commits, _ := repository.Log(&git.LogOptions{FileName: &fileName})
	// commits.ForEach(func(c *object.Commit) error {
	// 	f, err := c.Files()

	// 	if err != nil {
	// 		fmt.Println("Could not get files for commit: ", c.Hash)
	// 	} else {
	// 		f.ForEach(func(f *object.File) error {
	// 			if f.Name == fileName {
	// 				fileContents, _ := f.Contents()
	// 				render([]byte(fileContents), buffer)

	// 				fmt.Println("---", c.Hash.String())
	// 				os.MkdirAll(outputFolder+"/"+title+"/"+c.Hash.String()[0:8], os.ModePerm)
	// 				os.WriteFile(outputFolder+"/"+title+"/"+c.Hash.String()[0:8]+"/index.html", buffer.Bytes(), os.ModePerm)

	// 				buffer.Reset()
	// 			}
	// 			return nil
	// 		})
	// 	}

	// 	return nil
	// })

	contents, _ := os.ReadFile(articlePath)
	item := Article{
		ID:           uuid.NewV5(uuid.NamespaceURL, articlePath).String(),
		Path:         filePath,
		Title:        title,
		Folder:       dir,
		Size:         f.Size(),
		Type:         "article",
		FileModified: f.ModTime().UTC().Format(time.RFC3339),
		Source:       string(contents),
		Html:         "",
	}
	html := render(contents, item)

	os.MkdirAll(outputFolder+"/"+title, os.ModePerm)
	os.WriteFile(outputFolder+"/"+title+"/index.html", []byte(html), os.ModePerm)

	jsonData, _ := json_marshal(item)
	os.WriteFile(outputFolder+"/"+title+"/index.json", jsonData, os.ModePerm)

	messages <- item
}

func main() {
	articleRoot := flag.String("a", "", "Article root")
	outputFolder := flag.String("o", "", "Output folder")
	flag.Parse()

	if *articleRoot == "" {
		fmt.Println("You must give me an article folder")
		os.Exit(1)
	}

	if *outputFolder == "" {
		fmt.Println("You must give me an output folder")
		os.Exit(2)
	}

	// Open git repo
	r, _ := git.PlainOpen(*articleRoot)

	// Process all articles
	_, err := process(*articleRoot, *outputFolder, ".md", r)
	if err != nil {
		fmt.Println(err)
	}

	// Copy static assets
	fmt.Println("Copying assets...")
	copyErr := cp.Copy(*articleRoot+"/__assets", *outputFolder+"/assets")
	if copyErr != nil {
		fmt.Println("Could not copy assets: ", copyErr)
	}
}
