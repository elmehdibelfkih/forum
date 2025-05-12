package handler

import (
	"forum/internal/db"
	repo "forum/internal/repository"
	"forum/internal/utils"
	"net/http"
	"strconv"
)

// TODO: check the bad request example: change the post id in the heden form
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
	if !utils.ValidComment(r.FormValue("comment")) {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	postId, err := strconv.ParseInt(r.FormValue("post_id"), 10, 0)
	IsPostExist, err2 := db.IsPostExist(int(postId))
	if err != nil || !IsPostExist{
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}
	if err2 != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)

		return
	}
	err = db.AddNewComment(r.Context().Value(repo.USER_ID_KEY).(int), int(postId), r.FormValue("comment"))
	if err != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
