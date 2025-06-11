package handler

import (
	"fmt"
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

	// Filter query by id  ---> !!!
	Idpost := r.URL.Query().Get("Id")
	Id , err := strconv.Atoi(Idpost)
	if err != nil {
		fmt.Fprintf(w,"bad request")
		return
	}
	// Get post by id - with limit commenter ---> !!!
	post , err := db.GetPostByID(Id, userID)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	// her i add the data  
	confMap["Post"] = post

	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "post.html", confMap)

}
