package handler

import (
	db "forum/internal/db"
	errTmp "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) { // todo: check the methode
	if r.URL.Path != "/" {
		errTmp.TempErr(w, nil, http.StatusNotFound)
		return
	}
	data, err := db.GetAllPostsInfo()
	if err != nil {
		errTmp.TempErr(w, nil, http.StatusInternalServerError)
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false, "Posts": data})
		return
	}

	user_id, exist, err := db.SelectUserSession(sessionCookie.Value)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
		return
	}

	if !exist {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false, "Posts": data})
		return
	}

	user, err := db.GetUserInfo(user_id)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
		return
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": true, "Username": user.Username, "Posts": data})
}
