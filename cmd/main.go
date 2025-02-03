package main

import (
	running "github.com/bllooop/musiclibrary/internal/server"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

// @title MusicLibrary API
// @version 1.0
// @description API сервис онлайн библиотека песен

// @host localhost:8000
// @BasePath /

func main() {
	running.Run()
}
