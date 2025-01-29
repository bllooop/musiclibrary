package usecase

import (
	"errors"

	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
)

type MusicUsecase struct {
	repo repository.SongsLibrary
}

func NewMusicService(repo repository.SongsLibrary) *MusicUsecase {
	return &MusicUsecase{
		repo: repo,
	}
}

func (s *MusicUsecase) GetSongsLibrary(order, sort string, page int) (map[string]interface{}, error) {
	return s.repo.GetSongsLibrary(order, sort, page)
}

func (s *MusicUsecase) DeleteSong(songid int) error {
	return s.repo.DeleteSong(songid)
}
func (s *MusicUsecase) Update(songid int, input domain.UpdateSong) error {
	if input.Name == nil && input.Group == nil && input.ReleaseDate == nil && input.Text == nil && input.Link == nil {
		return errors.New("update params have no values")
	}
	return s.repo.Update(songid, input)
}
