package handler

import (
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetPostHandler(w, r)
	} else if r.Method == http.MethodPost {
		PostPostHandler(w, r)
	} else {
		forumerror.MethodNotAllowed(w, r)

		return
	}
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	var errMap map[string]any
	if r.Context().Value(repo.ERROR_CASE) != nil {
		errMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "new_post.html", errMap) // when u excute 2 template the get concatinated one in top of the other
}

func PostPostHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	// TODO: check if the inputs is valid
	// TODO: add the categories

	err := db.AddNewPost(user_id, r.FormValue("title"), r.FormValue("content"))
	if err != nil {
		forumerror.InternalServerError(w,r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
