package handler

import (
	"net/http"
	repo "forum/internal/repository"
)


func PostHandler(w *http.ResponseWriter, r *http.Request) {

	// <--> Get User id 
	userId := r.Context().Value(repo.USER_ID_KEY).(int)
	

}
