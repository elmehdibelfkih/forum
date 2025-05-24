package handler

import (
	"forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"forum/internal/utils"
	"net/http"
	"strconv"
)

// TODO: check the bad request example: change the post id in the heden form
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		forumerror.MethodNotAllowed(w, r)
		return
	}
	if !utils.ValidComment(r.FormValue("comment")) {
		forumerror.BadRequest(w, r)
		return
	}

	postId, err := strconv.ParseInt(r.FormValue("post_id"), 10, 0)
	IsPostExist, err2 := db.IsPostExist(int(postId))
	if err != nil || !IsPostExist {
		forumerror.BadRequest(w, r)
		return
	}
	if err2 != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if !utils.ValidComment(r.FormValue("comment")) {
		forumerror.BadRequest(w, r)
		return
	}
	err = db.AddNewComment(r.Context().Value(repo.USER_ID_KEY).(int), int(postId), r.FormValue("comment"))
	if err != nil {
		forumerror.InternalServerError(w,r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
