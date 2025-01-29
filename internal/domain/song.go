package domain

type Song struct {
	Id          int    `json:"-" db:"id"`
	Group       string `json:"group" binding:"required"`
	Name        string `json:"name" binding:"required"`
	ReleaseDate string `json:"date" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}

type UpdateSong struct {
	Group       *string `json:"group"`
	Name        *string `json:"name"`
	ReleaseDate *string `json:"date"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
}

type SongDetail struct {
	ReleaseDate string `json:"date" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link" binding:"required"`
}
