package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)


func main() {
	var pathDB string
	var fileMigrationPath string

	flag.StringVar(&pathDB, "storage-path", "", "table creation path")
	flag.StringVar(&fileMigrationPath, "migrations-path", "", "path to migration file")
	flag.Parse()

	if pathDB == "" || fileMigrationPath == "" {
		panic("the path of the file with migrations or the path for database creation is not specified")
	}

	migrateDb, err := migrate.New("file://"+fileMigrationPath, pathDB)
	if err != nil {
		panic(err)
	}

	if err := migrateDb.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	log.Println("migrations have been successfully applied")
}
