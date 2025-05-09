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
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false})
		return
	}

	user_id, exist, err := db.SelectUserSession(sessionCookie.Value)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
	}

	if !exist {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false})
		return
	}

	user, err := db.GetUserInfo(user_id)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
	}

	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": true, "Username": user.Username})
}
