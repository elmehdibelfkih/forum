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
	UserId     int
	Title      string
	Content    string
	Publisher  string
	Catigories []string
	Likes      int
	Dislikes   int
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
	INTERNAL_SERVER_ERROR_LOG_PATH = "./logs/internal_errors.log"

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
	// user not login or login !
	PUBLIC contextKey = "Public"


	INIT_FIELDS_QUERY = `INSERT OR IGNORE INTO categories (name) VALUES (?)`
	
)

// IT major fields
var IT_MAJOR_FIELDS = map[string]bool{
	"Software Engineering":                         true,
	"Artificial Intelligence and Machine Learning": true,
	"Data Science & Big Data":                      true,
	"Cybersecurity":                                true,
	"Networking & Telecommunications":              true,
	"Cloud Computing & Virtualization":             true,
	"DevOps & SRE":                                 true,
	"Database Systems":                             true,
	"Systems Programming":                          true,
	"Reverse Engineering":                          true,
	"Mobile & Embedded Development":                true,
	"IoT (Internet of Things)":                     true,
}

