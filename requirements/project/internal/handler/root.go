package handler

import (
	"errors"
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
	"strconv"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // FIXME: debagging
		// forumerror.NotFoundError(w, r)
		forumerror.InternalServerError(w, r, errors.New("test for runing in container"))
		return
	}
	if r.Method != http.MethodGet {
		forumerror.MethodNotAllowed(w, r)
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
	
	if r.URL.Query().Get("filter") == "All categories" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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
	filter := query.Get("filter")

	if filter != "" && filter != "Owned" && filter != "Likes" && !repo.IT_MAJOR_FIELDS[filter] {
		forumerror.BadRequest(w, r)
		return -1, errors.New("invalid filter")
	}

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

	count, err := db.GetPostsCount(query.Get("filter"), r.Context().Value(repo.USER_ID_KEY).(int))
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return -1, err
	}
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
	userId := r.Context().Value(repo.USER_ID_KEY).(int)

	filter := query.Get("filter")
	if filter == "" {
		confMap["Filter"] = "All Posts"
		data, err := db.GetAllPostsInfo(page, userId)
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
	confMap["Filter"] = filter
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
		data, err := db.GePostbycategory(filter, page, userId)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return err
		}
		confMap["Posts"] = data
	}

	return nil
}
