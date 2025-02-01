package usecase

import (
	"errors"

	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/bllooop/musiclibrary/internal/repository"
	logger "github.com/bllooop/musiclibrary/pkg/logging"
)

type MusicUsecase struct {
	repo repository.SongsLibrary
}

func NewMusicService(repo repository.SongsLibrary) *MusicUsecase {
	return &MusicUsecase{
		repo: repo,
	}
}

func (s *MusicUsecase) GetSongsLibrary(order, sort string, page int, name, group, text, releasedate, link string) (map[string]interface{}, error) {

	validSortFields := []string{"name", "artist", "releasedate"}
	validOrders := []string{"ASC", "DESC"}
	if !contains(validSortFields, sort) {
		return nil, errors.New("invalid sort value")
	}
	if !contains(validOrders, order) {
		return nil, errors.New("invalid order value")
	}
	logger.Log.Debug().Msgf("Successfully validated Sort: %s, Order: %s", sort, order)

	return s.repo.GetSongsLibrary(order, sort, page, name, group, text, releasedate, link)
}

func (s *MusicUsecase) DeleteSong(songid int) error {
	return s.repo.DeleteSong(songid)
}
func (s *MusicUsecase) Update(songid int, input domain.UpdateSong) error {
	if input.Name == nil && input.Group == nil && input.ReleaseDate == nil && input.Text == nil && input.Link == nil {
		return errors.New("update params have no values")
	}
	logger.Log.Debug().Msgf("Successfully validated Name: %s, Group: %s, Text: %s, Release Date: %s, Link: %s ", input.Name, input.Group, input.ReleaseDate, input.Text, input.Link)

	return s.repo.Update(songid, input)
}

func (s *MusicUsecase) CreateSong(song domain.UpdateSong, songDetail domain.UpdateSong) (int, error) {
	return s.repo.CreateSong(song, songDetail)
}

func (s *MusicUsecase) GetSongsById(songName string, begin, end int) ([]domain.Song, error) {
	return s.repo.GetSongsById(songName, begin, end)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
