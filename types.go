package main

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
