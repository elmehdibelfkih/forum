package handler

import (
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"html"
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
	var confMap = make(map[string]any)
	if r.Context().Value(repo.ERROR_CASE) != nil {
		confMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}
	confMap["Fields"] = repo.IT_MAJOR_FIELDS
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	user, err := db.GetUserInfo(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	confMap["Username"] = user.Username
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "new_post.html", confMap) // when u excute 2 template the get concatinated one in top of the other
}

func PostPostHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	// TODO: check if the inputs is valid
	// TODO: add the categories
	escapedTitle := html.EscapeString(r.FormValue("title"))
	escapedContent := html.EscapeString(r.FormValue("content"))
	postId, err := db.AddNewPost(user_id, escapedTitle, escapedContent)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	categories := r.Form["Categories"]

	for _, c := range categories {
		if !repo.IT_MAJOR_FIELDS[c] {
			forumerror.BadRequest(w, r)
			return
		}
	}

	err = db.MapPostWithCategories(postId, categories)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
