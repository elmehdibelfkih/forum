package repository

import (
	"database/sql"
	"html/template"
	"regexp"
)

var (
	EmailExp        *regexp.Regexp
	UsernameExp     *regexp.Regexp
	GLOBAL_TEMPLATE *template.Template
	DB              *sql.DB
)

type User struct {
	Id            int
	Username      string
	Email         string
	Password_hash string
	Created_at    string
	Updated_at    string
}

type Post struct {
	Id         int
	Title      string
	Content    string
	Publisher  string
	Catigory   string
	Likes      int
	Deslikes   int
	Comments   []map[string]string
	Created_at string
	Updated_at string
	IsEdited   bool
}

type PageData struct {
	Posts []Post
}

type contextKey string

const (
	DATABASE_NAME = "sqlite3"

	// database paths
	DATABASE_LOCATION        = "./database/forum.db"
	DATABASE_SCHEMA_LOCATION = "./database/schema.sql"

	// error messages
	FAILED_OPEN_DATABES     = "failed to open the database: %v"
	FAILED_CREAT_TABELS     = "failed to create tables: %v"
	FAILED_CLOSING_DATABASE = "closing database: %v"

	// templates paths
	TEMPLATE_PATHS = "./templates/*.html"

	// server helpers
	PORT               = ":8080"
	SERVER_RUN_MESSAGE = "\033[2mServer running on http://localhost:8080\033[0m"

	// context keys
	USER_ID_KEY contextKey = "user_id"
	ERROR_CASE  contextKey = "error_case"
)
