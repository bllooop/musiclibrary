package models

type Song struct {
	Id          int    `json:"-" db:"id"`
	Group       string `json:"group" binding:"required"`
	ReleaseDate string `json:"date" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}
