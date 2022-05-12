package main

import "time"

type Hierarchy struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type Article struct {
	ID           string      `json:"id"`
	URI          string      `json:"uri"`
	Size         int64       `json:"sizeInBytes"`
	Title        string      `json:"title"`
	FileModified time.Time   `json:"modified"`
	Source       string      `json:"source"`
	Html         string      `json:"html"`
	Hierarchy    []Hierarchy `json:"hierarchy"`
}

type Children struct {
	Articles []FolderThing `json:"articles"`
	Folders  []FolderThing `json:"folders"`
}

type Folder struct {
	ID        string      `json:"id"`
	URI       string      `json:"uri"`
	Title     string      `json:"title"`
	Hierarchy []Hierarchy `json:"hierarchy"`
	Children  Children    `json:"children"`
	README    string      `json:"readme"`
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

type Statistics struct {
	GenerationTime time.Duration `json:"generationTime"`
	ArticleCount   int           `json:"articleCount"`
	CPUCount       int           `json:"cpuCount"`
	MemoryInGB     int           `json:"memoryInGB"`
	Platform       string        `json:"platform"`
	Architecture   string        `json:"architecture"`
}

type BockConfig struct {
	articleRoot  string
	outputFolder string
	statistics   Statistics
}

type FolderThing struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	URI  string `json:"uri"`
}
