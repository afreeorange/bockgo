package main

import "time"

type Article struct {
	ID           string    `json:"id"`
	URI          string    `json:"uri"`
	Size         int64     `json:"sizeInBytes"`
	Title        string    `json:"title"`
	FileModified time.Time `json:"modified"`
	Source       string    `json:"source"`
	Html         string    `json:"html"`
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
	GenerationTime time.Duration
	ArticleCount   int
	CPUs           int
	Memory         int
	Platform       string
	Architecture   string
}
