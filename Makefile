build:
		docker compose build musiclibrary
run:
		docker compose up musiclibrary
postgrescont:
		docker run --name=db -e POSTGRES_PASSWORD='54321' -p 5436:5432 -d postgres
migrate:	
		goose -dir migrations postgres "postgres://postgres:54321@localhost:5436/postgres?sslmode=disable" up
migrate-down:	
		goose -dir migrations postgres "postgres://postgres:54321@localhost:5436/postgres?sslmode=disable" down