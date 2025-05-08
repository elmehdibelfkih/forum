package errors

import (
	"log"
	"net/http"
)

func TempErr(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}
