package usecase

import (
	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
)

type SongsLibrary interface {
	GetSongsLibrary(order, sortby string, page int, name, group, text, releasedate, link string) (map[string]interface{}, error)
	DeleteSong(songid int) error
	Update(songid int, input domain.UpdateSong) error
	CreateSong(song domain.UpdateSong, songDetail domain.UpdateSong) (int, error)
	GetSongsById(songName string, begin, end int) ([]domain.Verses, error)
}

type Usecase struct {
	SongsLibrary
}

func NewService(repository *repository.Repository) *Usecase {
	return &Usecase{
		SongsLibrary: NewMusicService(repository),
	}
}
