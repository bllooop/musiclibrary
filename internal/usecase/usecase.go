package usecase

import (
	models "github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
)

type GetSongs interface {
	GetSongsLibrary(order, sortby string, page int) ([]models.Song, error)
}

type Usecase struct {
	GetSongs
}

func NewService(repository *repository.Repository) *Usecase {
	return &Usecase{
		GetSongs: NewMusicService(repository),
	}
}
