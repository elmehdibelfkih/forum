package handler

import (
	forumerror "forum/internal/error"
	"net/http"
	"os"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		forumerror.MethodNotAllowed(w, r)
	}
	url := r.URL.Path[1:]
	file, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			forumerror.NotFoundError(w, r)
		}
		forumerror.InternalServerError(w, r, err)
	}
	if file.IsDir() {
		forumerror.NotFoundError(w, r)
	}
	http.ServeFile(w, r, url)
}
