package main

import "time"

type Hierarchy struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}

type Revised struct {
	Modified time.Time `json:"modified"`
	Created  time.Time `json:"created"`
}

type Revision struct {
	Id      string `json:"id"`
	ShortId string `json:"shortId"`
	Date    string `json:"date"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Article struct {
	FileModified time.Time   `json:"modified"`
	FileCreated  time.Time   `json:"modified"`
	Hierarchy    []Hierarchy `json:"hierarchy"`
	Html         string      `json:"html"`
	ID           string      `json:"id"`
	Size         int64       `json:"sizeInBytes"`
	Source       string      `json:"source"`
	Title        string      `json:"title"`
	URI          string      `json:"uri"`
	Revisions    []Revision  `json:"revisions"`
}

type Children struct {
	Articles []FolderThing `json:"articles"`
	Folders  []FolderThing `json:"folders"`
}

type Folder struct {
	Children  Children    `json:"children"`
	Hierarchy []Hierarchy `json:"hierarchy"`
	ID        string      `json:"id"`
	README    string      `json:"readme"`
	Title     string      `json:"title"`
	URI       string      `json:"uri"`
}

type Statistics struct {
	Architecture   string        `json:"architecture"`
	ArticleCount   int           `json:"articleCount"`
	BuildDate      time.Time     `json:"buildTime"`
	CPUCount       int           `json:"cpuCount"`
	GenerationTime time.Duration `json:"generationTime"`
	MemoryInGB     int           `json:"memoryInGB"`
	Platform       string        `json:"platform"`
}

type BockConfig struct {
	articleRoot  string
	outputFolder string
	makeJSON     bool
	statistics   Statistics
	started      time.Time
}

type FolderThing struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}
