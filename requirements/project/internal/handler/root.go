package handler

import (
	"errors"
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
	if r.Context().Value(repo.USER_ID_KEY).(int) == -1 {
		confMap["Authenticated"] = false
	} else {
		confMap["Authenticated"] = true
		confMap["Username"] = r.Context().Value(repo.USER_NAME).(string)
	}

	page, err := Pagination(w, r, confMap)
	if err != nil {
		return
	}

	err = GetPostsByFilter(w, r, confMap, page)
	if err != nil {
		return
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", confMap)
}

func Pagination(w http.ResponseWriter, r *http.Request, confMap map[string]any) (int, error) {
	query := r.URL.Query()
	page := 1

	if pageStr := query.Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			forumerror.BadRequest(w, r)
			return -1, err
		}
		page = p
	}

	confMap["CurrentPage"] = page
	confMap["PrintCurrentPage"] = page != 1
	confMap["HasPrev"] = page > 1

	if page > 1 {
		prevQuery := r.URL.Query()
		prevQuery.Set("page", strconv.Itoa(page-1))
		confMap["PrevPage"] = r.URL.Path + "?" + prevQuery.Encode()
	}

	count, err := db.GetPostsCount(query.Get("filter"))

	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return -1, err
	}
	// todo: fix the count logic
	hasNext := count > page*repo.PAGE_POSTS_QUANTITY
	confMap["HasNext"] = hasNext
	if hasNext {
		nextQuery := r.URL.Query()
		nextQuery.Set("page", strconv.Itoa(page+1))
		confMap["NextPage"] = r.URL.Path + "?" + nextQuery.Encode()
	}

	return page, nil
}

func GetPostsByFilter(w http.ResponseWriter, r *http.Request, confMap map[string]any, page int) error {
	query := r.URL.Query()

	filter := query.Get("filter")
	if filter == "" {
		data, err := db.GetAllPostsInfo(page)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return err
		}
		confMap["Posts"] = data
		return nil
	}
	if filter != "Owned" && filter != "Likes" && !repo.IT_MAJOR_FIELDS[filter] {
		forumerror.BadRequest(w, r)
		return errors.New("resource not found")
	}
	userId := r.Context().Value(repo.USER_ID_KEY).(int)
	switch filter {
	case "Owned":
		if !confMap["Authenticated"].(bool) {
			forumerror.Unauthorized(w, r)
			return errors.New("err")
		}
		data, err := db.Getpostbyowner(userId, page)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return err
		}
		confMap["Posts"] = data
	case "Likes":
		if !confMap["Authenticated"].(bool) {
			forumerror.Unauthorized(w, r)
			return errors.New("err")
		}
		data, err := db.Getposbytlikes(userId, page)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return err
		}
		confMap["Posts"] = data
	default:
		data, err := db.GePostbycategory(filter, page)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return err
		}
		confMap["Posts"] = data
	}

	return nil
}
