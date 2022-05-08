package main

type Article struct {
	ID           string `json:"id"`
	URI          string `json:"uri"`
	Size         int64  `json:"sizeInBytes"`
	Title        string `json:"title"`
	FileModified string `json:"modified"`
	Source       string `json:"source"`
	Html         string `json:"html"`
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
