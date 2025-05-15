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
)

// func NotFoundError(w http.ResponseWriter) {
// 	notFoundError := Error{
// 		StatusCode:       404,
// 		StatusText:       "Not Found",
// 		ErrorMessage:     "The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.",
// 		ErrorTitle:       "Oops! Page Not Found",
// 		ErrorDescription: "We couldn't find the page you were looking for. Please check the URL for any mistakes or go back to the homepage.",
// 	}
// 	w.WriteHeader(notFoundError.StatusCode)
// 	ERRORTMPL.Execute(w, notFoundError)
// }

// func MethodNotAllowed(w http.ResponseWriter) {
// 	methodNotAllowed := Error{
// 		StatusCode:       405,
// 		StatusText:       "Method not allowed",
// 		ErrorMessage:     "The HTTP method used for this request is not allowed on this resource. Please use a different method.",
// 		ErrorTitle:       "Oops! method not allowed",
// 		ErrorDescription: "Please send a POST request to this endpoint with the required data for processing. Other HTTP methods are not allowed.",
// 	}
// 	w.WriteHeader(methodNotAllowed.StatusCode)
// 	ERRORTMPL.Execute(w, methodNotAllowed)
// }

// func InternalServerError(w http.ResponseWriter) {
// 	internalServerError := Error{
// 		StatusCode:       500,
// 		StatusText:       "Internal Server Error",
// 		ErrorMessage:     "An unexpected error occurred while processing your request.",
// 		ErrorTitle:       "Oops! Internal Server Error",
// 		ErrorDescription: "Something went wrong on our end. Please try again later or contact support if the issue persists.",
// 	}
// 	w.WriteHeader(internalServerError.StatusCode)
// 	ERRORTMPL.Execute(w, internalServerError)
// }

// func BadRequest(w http.ResponseWriter) {
// 	badRequest := Error{
// 		StatusCode:       400,
// 		StatusText:       "Bad Request",
// 		ErrorMessage:     "The Server recieved a Bad request!!",
// 		ErrorTitle:       "Oops! Bad Request",
// 		ErrorDescription: "Please make sure you respect the input limits",
// 	}
// 	w.WriteHeader(badRequest.StatusCode)
// 	ERRORTMPL.Execute(w, badRequest)
// }