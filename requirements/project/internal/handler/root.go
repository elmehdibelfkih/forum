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
	data := repo.PageData{
		Posts: []repo.Post{
			{
				Id:        1,
				Title:     "Understanding Go Templates",
				Content:   "Templates in Go let you separate logic and HTML... Templates in Go let you separate logic and HTML... Templates in Go let you separate logic and HTML... Templates in Go let you separate logic and HTML... ",
				Publisher: "El Mehdi",
				Catigories:  []string{"Programming"},
				Likes:     42,
				Deslikes:  1,
				Comments: []map[string]string{
					{"hamid": "Great post!"},
					{"3li": "Thanks for sharing."},
				},
				Created_at: "2025-05-11",
				Updated_at: "2025-05-11",
				IsEdited:   false,
			},
			{
				Id:        2,
				Title:     "test",
				Content:   "kanjarbo wach hadchi khdam wla la lakan khdam rah nadi hadchi ",
				Publisher: "El chapo",
				Catigories:  []string{"walo"},
				Likes:     37,
				Deslikes:  13,
				Comments: []map[string]string{
					{"hamid": "Great post!"},
					{"3li": "Thanks for sharing."},
				},
				Created_at: "2025-05-01",
				Updated_at: "2025-05-13",
				IsEdited:   true,
			},
		},
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false, "Posts": data})
		return
	}

	user_id, exist, err := db.SelectUserSession(sessionCookie.Value)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
	}

	if !exist {
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": false, "Posts": data})
		return
	}

	user, err := db.GetUserInfo(user_id)

	if err != nil {
		errTmp.TempErr(w, err, http.StatusInternalServerError)
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": true, "Username": user.Username, "Posts": data})
}
