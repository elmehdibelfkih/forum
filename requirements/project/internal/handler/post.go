package handler

import (
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		forumerror.MethodNotAllowed(w, r)
		return
	}

	var confMap = make(map[string]any)

	if r.Context().Value(repo.USER_ID_KEY).(int) == -1 {
		confMap["Authenticated"] = false
	} else {
		confMap["Authenticated"] = true
		confMap["Username"] = r.Context().Value(repo.USER_NAME).(string)
	}


	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "post.html", confMap)

}
