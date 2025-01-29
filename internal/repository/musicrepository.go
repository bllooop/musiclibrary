package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bllooop/musiclibrary/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (r *MusicPostgres) GetSongsLibrary(order, sort string, page int) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	limit := 10

	offset := limit * (page - 1)
	data["Page"] = r.pagination("users", limit, page)
	var songs []domain.Song
	query := ""
	query = fmt.Sprintf(`SELECT name, artist, releasedate, releasedate, link, text
		FROM %s ORDER BY %s %s limit %d offset %d`, songsListTable, sort, order, limit, offset)
	row, err := r.pg.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var err error
		k := domain.Song{}
		err = row.Scan(&k.Id, &k.Name, &k.Group, &k.ReleaseDate, &k.Link, &k.Text)
		if err != nil {
			return nil, err
		}
		songs = append(songs, k)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	data["Songs"] = songs
	return data, nil
}
func (r *MusicPostgres) Update(songid int, input domain.UpdateSong) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *&input.Name)
		argId++
	}
	if input.Group != nil {
		setValues = append(setValues, fmt.Sprintf("artist=$%d", argId))
		args = append(args, *input.Group)
		argId++
	}
	if input.ReleaseDate != nil {
		setValues = append(setValues, fmt.Sprintf("releasedate=$%d", argId))
		args = append(args, *input.ReleaseDate)
		argId++
	}
	if input.Text != nil {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, *input.Text)
		argId++
	}
	if input.Link != nil {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argId))
		args = append(args, *input.Link)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", songsListTable, setQuery, argId)
	args = append(args, songid)
	_, err := r.pg.Exec(context.Background(), query, args...)
	return err
}

func (r *MusicPostgres) DeleteSong(songid int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", songsListTable)
	_, err := r.pg.Exec(context.Background(), query, songid)
	return err
}

func (r *MusicPostgres) pagination(table string, limit, page int) *Pagination {
	var (
		tmpl        = Pagination{}
		recordcount int
	)

	sqltable := fmt.Sprintf("SELECT count(id) FROM %s", table)

	r.pg.QueryRow(context.Background(), sqltable).Scan(&recordcount)

	total := (recordcount / limit)

	remainder := (recordcount % limit)
	if remainder == 0 {
		tmpl.TotalPage = total
	} else {
		tmpl.TotalPage = total + 1
	}

	tmpl.CurrentPage = page
	tmpl.RecordPerPage = limit

	if page <= 0 {
		tmpl.Next = page + 1
	} else if page < tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = page + 1
	} else if page == tmpl.TotalPage {
		tmpl.Previous = page - 1
		tmpl.Next = 0
	}

	return &tmpl
}
