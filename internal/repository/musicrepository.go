package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Pagination struct {
	Next          int
	Previous      int
	RecordPerPage int
	CurrentPage   int
	TotalPage     int
}
type MusicPostgres struct {
	pg *pgxpool.Pool
}

func NewMusicPostgres(pg *pgxpool.Pool) *MusicPostgres {
	return &MusicPostgres{
		pg: pg,
	}
}

func (r *MusicPostgres) GetSongsLibrary(order, sortby string, page int)
