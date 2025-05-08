package db

import (
	"database/sql"
	repo "forum/internal/repository"
	"log"
	"os"
)

func InitDB(datasource string) {
	var err error
	repo.DB, err = sql.Open(repo.DATABASE_NAME, datasource)
	if err != nil {
		log.Fatalf(repo.FAILED_OPEN_DATABES, err)
	}
	err = CreateTable(repo.DB)
	if err != nil {
		log.Fatalf(repo.FAILED_CREAT_TABELS, err)
	}
}

func CreateTable(db *sql.DB) error {
	schema, err := os.ReadFile(repo.DATABASE_SCHEMA_LOCATION)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}

func CloseDB() {
	err := repo.DB.Close()
	if err != nil {
		log.Fatalf(repo.FAILED_CLOSING_DATABASE, err)
	}
}
