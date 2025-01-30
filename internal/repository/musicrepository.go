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

func (r *MusicPostgres) CreateSong(song domain.UpdateSong, songDetail domain.UpdateSong) (int, error) {
	tr, err := r.pg.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (name, artist, releasedate, text, link) VALUES ($1,$2,$3,$4,$5) RETURNING *", songsListTable)
	row := tr.QueryRow(context.Background(), createListQuery, song.Name, song.Group, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
	if err := row.Scan(&id); err != nil {
		tr.Rollback(context.Background())
		return 0, err
	}
	return id, tr.Commit(context.Background())
}
func (r *MusicPostgres) GetSongsLibrary(order, sort string, page int, name, group, text, releasedate, link string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	limit := 10
	offset := limit * (page - 1)
	data["Page"] = r.pagination("users", limit, page)
	var filters []string
	var args []interface{}
	query := ""
	query = fmt.Sprintf(`SELECT name, artist, releasedate, link, text FROM %s`, songsListTable)

	if name != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE $%d", len(args)+1))
		args = append(args, "%"+name+"%")
	}
	if group != "" {
		filters = append(filters, fmt.Sprintf("artist ILIKE $%d", len(args)+1))
		args = append(args, "%"+group+"%")
	}
	if text != "" {
		filters = append(filters, fmt.Sprintf("text = $%d", len(args)+1))
		args = append(args, text)
	}
	if releasedate != "" {
		filters = append(filters, fmt.Sprintf("releasedate = $%d", len(args)+1))
		args = append(args, releasedate)
	}
	if link != "" {
		filters = append(filters, fmt.Sprintf("link = $%d", len(args)+1))
		args = append(args, link)
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	query += fmt.Sprintf(` ORDER BY %s %s limit %d offset %d`, sort, order, limit, offset)

	var songs []domain.Song
	row, err := r.pg.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		k := domain.Song{}
		err := row.Scan(&k.Id, &k.Name, &k.Group, &k.ReleaseDate, &k.Link, &k.Text)
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
		args = append(args, *input.Name)
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
