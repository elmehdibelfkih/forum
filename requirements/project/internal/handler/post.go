package handler

import (
	forumerror "forum/internal/error"
	db "forum/internal/db"
	repo "forum/internal/repository"
	"net/http"
	"strconv"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		forumerror.MethodNotAllowed(w, r)
		return
	}

	var confMap = make(map[string]any)
	userID := r.Context().Value(repo.USER_ID_KEY).(int) 
	if  r.Context().Value(repo.USER_ID_KEY).(int) == -1 {
		confMap["Authenticated"] = false
	} else {
		confMap["Authenticated"] = true
		confMap["Username"] = r.Context().Value(repo.USER_NAME).(string)

	}

	// Filter Query By Id  ---> !!!
	Idpost := r.URL.Query().Get("Id")
	Id , err := strconv.Atoi(Idpost)
	if err != nil {
		forumerror.BadRequest(w,r)
		return
	}
	// Get post by id - separate with comment !!!
	post , err := db.GetPostByID(Id, userID)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	// Her I Add The Data  
	confMap["Post"] = post

	// her i will handle the pagination 10 by 10 !!
	// Handle pagination: get page from query param, default to 1
    page := 1
    if pageStr := r.URL.Query().Get("page"); pageStr != "" {
        p, err := strconv.Atoi(pageStr)
        if err != nil || p < 1 {
            forumerror.BadRequest(w, r)
            return
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

    // Fetch paginated comments for the post (10 per page)
    comments, totalComments, err := db.GetCommentsByPostPaginated(Id, page, 10)
    if err != nil {
        forumerror.InternalServerError(w, r, err)
        return
    }
    confMap["Comments"] = comments

    // Determine if there is a next page
    hasNext := totalComments > page*10
    confMap["HasNext"] = hasNext

    if hasNext {
        nextQuery := r.URL.Query()
        nextQuery.Set("page", strconv.Itoa(page+1))
        confMap["NextPage"] = r.URL.Path + "?" + nextQuery.Encode()
    } 
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "post.html", confMap)
}
