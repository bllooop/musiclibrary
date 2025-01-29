package usecase

import (
	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
)

type SongsLibrary interface {
	GetSongsLibrary(order, sortby string, page int) (map[string]interface{}, error)
	DeleteSong(songid int) error
	Update(songid int, input domain.UpdateSong) error
}

type Usecase struct {
	SongsLibrary
}

func NewService(repository *repository.Repository) *Usecase {
	return &Usecase{
		SongsLibrary: NewMusicService(repository),
	}
}
