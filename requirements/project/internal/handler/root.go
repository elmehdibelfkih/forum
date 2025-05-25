package handler

import (
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
	"strconv"
)

// todo: fix the pagination logic
func RootHandler(w http.ResponseWriter, r *http.Request) { // todo: check the methode
	var confMap = make(map[string]any)

	// pagination
	query := r.URL.Query()
	pageStr := query.Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			forumerror.BadRequest(w, r)
		}
	}
	confMap["CurrentPage"] = pageStr
	if page == 1 {
		confMap["HasPrev"] = false
	} else {
		confMap["HasPrev"] = true
		confMap["PrevPage"] = strconv.Itoa(page-1)
	}
	confMap["HasNext"] = true
	confMap["NextPage"] = strconv.Itoa(page+1)
	// end of pagination

	confMap["Fields"] = repo.IT_MAJOR_FIELDS
	if r.URL.Path != "/" {
		forumerror.NotFoundError(w, r)
		return
	}
	data, err := db.GetAllPostsInfo(page)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	sessionCookie, err := r.Cookie("session_token")
	confMap["Authenticated"] = false
	confMap["Posts"] = data
	if err != nil || sessionCookie.Value == "" {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", confMap)
		return
	}
	user_id, exist, err := db.SelectUserSession(sessionCookie.Value)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if !exist {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", confMap)
		return
	}
	user, err := db.GetUserInfo(user_id)
	confMap["Username"] = user.Username
	confMap["Authenticated"] = true
	confMap["Posts"] = data
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", confMap)
}
