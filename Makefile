default: run
.PHONY: run

PATH_DB=postgres://evans:evans@localhost:8081/postgres?sslmode=disable
FILE_MIGRATIONS = ./migrations

run:
	go run cmd/main.go -config ./config.env

migrate:	# is to perform the migration at local startup
	go run cmd/migrator/migrate.go -storage-path $(PATH_DB) -migrations-path $(FILE_MIGRATIONS)

run-docker-compose:
	docker compose --env-file config.env up --build -d