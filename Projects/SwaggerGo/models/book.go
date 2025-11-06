package models

// Book is a model for book
type Book struct {
	Title     string `json:"title" example:"Book A"`
	Author    string `json:"author" example:"Dino"`
	Publisher string `json:"publisher" example:"Creative Company"`
}
