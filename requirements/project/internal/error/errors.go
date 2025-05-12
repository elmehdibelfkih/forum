package errors

import (
	"net/http"
)

func TempErr(w http.ResponseWriter, err error, code int) {
	// log.Println("error >>", err)
	http.Error(w, http.StatusText(code), code)
}
