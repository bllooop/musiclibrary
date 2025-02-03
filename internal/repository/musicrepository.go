package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bllooop/musiclibrary/internal/domain"
	logger "github.com/bllooop/musiclibrary/pkg/logging"
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
func (r *MusicPostgres) GetSongsById(songName string, begin, end int) ([]domain.Verses, error) {
	var verses []domain.Verses
	query := fmt.Sprintf(`WITH split_text AS (
    SELECT 
        verse, 
        row_number() OVER () AS verse_number
    FROM (
        SELECT unnest(string_to_array(text, '\n\n')) AS verse
        FROM %s
        WHERE name ILIKE $1
    ) AS unnested_lines
)
SELECT verse_number, verse
FROM split_text
WHERE verse_number BETWEEN $2 AND $3
ORDER BY verse_number;`, songsListTable)
	logger.Log.Debug().Str("query", query).Str("song_name", songName).
		Int("begin", begin).Int("end", end).Msg("Fetching song verses / Получение куплетов песни")
	logger.Log.Info().Msg("Executing query for get songs / Выполнение запроса на получение песен")
	rows, err := r.pg.Query(context.Background(), query, "%"+songName+"%", begin, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var k domain.Verses
		if err := rows.Scan(&k.Number, &k.Verse); err != nil {
			logger.Log.Error().Err(err).Msg("Error scanning row / Ошибка сканирования строки")
			return nil, err
		}
		verses = append(verses, k)
	}
	logger.Log.Debug().Int("songs_found", len(verses)).Msg("Successfully fetched songs by ID / Успешно найдены песни по идентификатору")
	return verses, nil
}

func (r *MusicPostgres) CreateSong(song domain.UpdateSong, songDetail domain.UpdateSong) (int, error) {
	tr, err := r.pg.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	var id int
	*songDetail.Text = strings.ReplaceAll(*songDetail.Text, "'", "''")
	createListQuery := fmt.Sprintf("INSERT INTO %s (name, artist, releasedate, text, link) VALUES ($1,$2,$3,$4,$5) RETURNING id", songsListTable)
	logger.Log.Debug().Str("query", createListQuery).Msg("Executing CreateSong query / Выполнение запроса CreateSong")
	row := tr.QueryRow(context.Background(), createListQuery, song.Name, song.Group, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
	if err := row.Scan(&id); err != nil {
		tr.Rollback(context.Background())
		return 0, err
	}
	err = tr.Commit(context.Background())
	if err != nil {
		return 0, err
	}
	logger.Log.Debug().Int("song_id", id).Msg("Successfully created song / Успешно сохраненная песня")
	return id, nil
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
	logger.Log.Debug().
		Str("query", query).
		Interface("args", args).
		Int("page", page).
		Msg("Executing GetSongsLibrary query")

	var songs []domain.Song
	row, err := r.pg.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		k := domain.Song{}
		err := row.Scan(&k.Name, &k.Group, &k.ReleaseDate, &k.Link, &k.Text)
		if err != nil {
			return nil, err
		}
		songs = append(songs, k)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	logger.Log.Debug().Int("songs_count", len(songs)).Msg("Successfully fetched songs from library / Успешно получены песни из библиотеки")
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
	logger.Log.Debug().
		Str("query", query).
		Msg("Executing Update query")
	args = append(args, songid)
	_, err := r.pg.Exec(context.Background(), query, args...)
	return err
}

func (r *MusicPostgres) DeleteSong(songid int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", songsListTable)
	logger.Log.Debug().
		Str("query", query).
		Msg("Executing Delete query")
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
