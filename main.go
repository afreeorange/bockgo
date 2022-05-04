package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	uuid "github.com/satori/go.uuid"
)

type Article struct {
	ID           string     `json:"id"`
	Folder       string     `json:"folder"`
	Path         []string   `json:"path"`
	Size         int64      `json:"sizeInBytes"`
	Title        string     `json:"title"`
	Type         string     `json:"type"`
	FileModified string     `json:"modified"`
	Source       string     `json:"source"`
	Html         string     `json:"html"`
	Revisions    []Revision `json:"revisions"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Revision struct {
	Id      string `json:"id"`
	ShortId string `json:"shortId"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Author  Author `json:"author"`
}

func processArticle(
	path string,
	articleRoot string,
	buffer *bytes.Buffer,
	f os.FileInfo,
	repository *git.Repository,
	messages chan Article,
) {
	dir := strings.Replace(filepath.Dir(path), articleRoot, "", -1)
	fileName := f.Name()
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	filePath := strings.Split(dir, "/")[1:]

	fmt.Println("Doing", fileName)

	commits, _ := repository.Log(&git.LogOptions{FileName: &fileName})
	commits.ForEach(func(c *object.Commit) error {
		f, err := c.Files()

		if err != nil {
			fmt.Println("Could not get files for commit: ", c.Hash)
		} else {
			f.ForEach(func(f *object.File) error {
				if f.Name == fileName {
					fileContents, _ := f.Contents()
					render([]byte(fileContents), buffer)

					fmt.Println("---", c.Hash.String())
					os.MkdirAll("/Users/nikhilanand/Desktop/temp/"+title+"/"+c.Hash.String()[0:8], os.ModePerm)
					os.WriteFile("/Users/nikhilanand/Desktop/temp/"+title+"/"+c.Hash.String()[0:8]+"/index.html", buffer.Bytes(), os.ModePerm)

					buffer.Reset()
				}
				return nil
			})
		}

		return nil
	})

	contents, _ := os.ReadFile(path)
	render(contents, buffer)
	item := Article{
		ID:           uuid.NewV5(uuid.NamespaceURL, path).String(),
		Path:         filePath,
		Title:        title,
		Folder:       dir,
		Size:         f.Size(),
		Type:         "article",
		FileModified: f.ModTime().UTC().Format(time.RFC3339),
		Source:       string(contents),
		Html:         buffer.String(),
	}

	os.MkdirAll("/Users/nikhilanand/Desktop/temp/"+title, os.ModePerm)
	os.WriteFile("/Users/nikhilanand/Desktop/temp/"+title+"/index.html", buffer.Bytes(), os.ModePerm)

	jsonData, _ := json_marshal(item)
	os.WriteFile("/Users/nikhilanand/Desktop/temp/"+title+"/index.json", jsonData, os.ModePerm)

	buffer.Reset()
	messages <- item
}

func glob(articleRoot string, extension string, r *git.Repository) ([]Article, error) {
	files := []Article{}
	var buffer bytes.Buffer

	messages := make(chan Article)

	err := filepath.Walk(articleRoot, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == extension {
			go processArticle(path, articleRoot, &buffer, f, r, messages)

			item := <-messages
			files = append(files, item)
		}

		return nil
	})

	return files, err
}

// func getRevisions(r *git.Repository, fileName string) {
// 	commits, _ := r.Log(&git.LogOptions{FileName: &fileName})

// 	commits.ForEach(func(c *object.Commit) error {
// 		f, err := c.Files()

// 		if err != nil {
// 			fmt.Println("Could not get files for commit: ", c.Hash)
// 		} else {
// 			f.ForEach(func(f *object.File) error {
// 				if f.Name == fileName {
// 					fmt.Println(c.Hash)
// 					fmt.Println(c.Committer.When)
// 					fmt.Println(f.Contents())
// 					fmt.Println("---")
// 				}
// 				return nil
// 			})
// 		}

// 		return nil
// 	})
// }

func main() {
	articleRoot := "/Users/nikhilanand/personal/wiki.nikhil.io.articles"
	r, _ := git.PlainOpen(articleRoot)

	list, err := glob(articleRoot, ".md", r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("I found %s articles in %s\n", fmt.Sprint(len(list)), articleRoot)
}
