build:
		docker compose build selling
run:
		docker compose up selling
postgrescont:
		docker run --name=db -e POSTGRES_PASSWORD='54321' -p 5436:5432 -d postgres
migrate:	
		migrate -path ./migrations -database 'postgres://postgres:54321@localhost:5436/postgres?sslmode=disable' up
migrate-down:	
		migrate -path ./migrations -database 'postgres://postgres:54321@localhost:5436/postgres?sslmode=disable' down