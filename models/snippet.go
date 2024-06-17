package models

type Snippet struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Content     string `json:"content"`
	Date        string `json:"dateTime"`
	Notes       string `json:"notes"`
}