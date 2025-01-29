package repository

import (
	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	SongsLibrary
}

type SongsLibrary interface {
	GetSongsLibrary(order, sort string, page int) (map[string]interface{}, error)
	DeleteSong(songid int) error
	Update(songid int, input domain.UpdateSong) error
}

func NewRepository(pg *pgxpool.Pool) *Repository {
	return &Repository{
		SongsLibrary: NewMusicPostgres(pg),
	}
}
