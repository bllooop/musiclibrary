package repository

import (
	models "github.com/bllooop/musiclibrary/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	GetSongs
}

type GetSongs interface {
	GetSongsLibrary(order, sort string, page int) ([]models.Song, error)
}

func NewRepository(pg *pgxpool.Pool) *Repository {
	return &Repository{
		GetSongs: NewMusicPostgres(pg),
	}
}
