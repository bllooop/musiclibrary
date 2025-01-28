package usecase

import (
	models "github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
)

type MusicUsecase struct {
	repo repository.GetSongs
}

func NewMusicService(repo repository.GetSongs) *MusicUsecase {
	return &MusicUsecase{
		repo: repo,
	}
}

func (s *MusicUsecase) GetSongsLibrary(order, sort string, page int) ([]models.Song, error) {
	return s.repo.GetSongs(order, sort, page)
}
