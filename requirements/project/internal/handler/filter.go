package handler

import (
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"forum/internal/utils"
	"net/http"
)

func Selectfilter(w http.ResponseWriter, r *http.Request) {
	// check methoud
	if r.Method != http.MethodGet {
		forumerror.MethodNotAllowed(w, r)
		return
	}
	// get auth & userid from context
	auth, _ := r.Context().Value(repo.PUBLIC).(bool)
	var userID int
	if auth {
		userID, _ = r.Context().Value(repo.USER_ID_KEY).(int)
	}

	// her is get filter in Query
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		forumerror.BadRequest(w, r)
		return
	}

	var posts repo.PageData
	var err error

	switch filter {
	case "Likes":
		if !auth {
			forumerror.Unauthorized(w, r)
			return
		}
		posts, err = db.Getposbytlikes(userID)

	case "Owned":
		if !auth {
			forumerror.Unauthorized(w, r)
			return
		}
		posts, err = db.Getpostbyowner(userID)

	default:
		if !utils.Contain(filter) {
			forumerror.BadRequest(w, r)
			return
		}
		posts, err = db.GePostbycategory(filter)
	}

	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	// if authenticated "login", fetch the name of user from database !
	username := ""
	if auth {
		var user repo.User
		user, err = db.GetUserInfo(userID)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return
		}
		username = user.Username
	}
	// her is add the data into map for template !
	data := map[string]any{
		"Authenticated": auth,
		"Username":      username,
		"Posts":         posts,
	}
	err = repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
}
