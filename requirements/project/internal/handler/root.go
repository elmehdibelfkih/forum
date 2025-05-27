package handler

import (
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
	"strconv"
)

func RootHandler(w http.ResponseWriter, r *http.Request) { // todo: check the methode
	if r.URL.Path != "/" {
		forumerror.NotFoundError(w, r)
		return
	}

	var confMap = make(map[string]any)

	confMap["Fields"] = repo.IT_MAJOR_FIELDS

	page, err := Pagination(w, r, confMap)
	if err != nil {
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

func Pagination(w http.ResponseWriter, r *http.Request, confMap map[string]any) (int, error) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			forumerror.BadRequest(w, r)
			return -1, err
		}
	}
	confMap["CurrentPage"] = pageStr
	confMap["PrintCurrentPage"] = page != 1
	if page == 1 {
		confMap["HasPrev"] = false
	} else {
		confMap["HasPrev"] = true
		confMap["PrevPage"] = strconv.Itoa(page - 1)
	}

	count, err := db.GetPostsCount()
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return -1, err

	}

	confMap["HasNext"] = count > page * repo.PAGE_POSTS_QUANTITY
	confMap["NextPage"] = strconv.Itoa(page + 1)
	return page, nil
}
