package repository

import (
	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	SongsLibrary
}

type SongsLibrary interface {
	GetSongsLibrary(order, sort string, page int, name, group, text, releasedate, link string) (map[string]interface{}, error)
	DeleteSong(songid int) error
	Update(songid int, input domain.UpdateSong) error
	CreateSong(song domain.UpdateSong, songDetail domain.UpdateSong) (int, error)
	GetSongsById(songName string, begin, end int) ([]domain.Verses, error)
}

func NewRepository(pg *pgxpool.Pool) *Repository {
	return &Repository{
		SongsLibrary: NewMusicPostgres(pg),
	}
}
