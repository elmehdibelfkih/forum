package auth

import (
	"context"
	"net/http"
	"time"

	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	utils "forum/internal/utils"
)

func SwitchLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ServLogin(w, r)
	case http.MethodPost:
		SubmitLogin(w, r)
	default:
		forumerror.MethodNotAllowed(w, r)
	}
}

func ServLogin(w http.ResponseWriter, r *http.Request) {
	var errMap map[string]any
	if r.Context().Value(repo.ERROR_CASE) != nil {
		errMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "login.html", errMap)
}

func SubmitLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	var userId int
	var hash string
	var err error
	exist, err := db.AlreadyExists(username, username)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	switch {
	case (!utils.ValidUsername(username) && !utils.ValidEmail(username)):
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "invalid usename or email plese try again"})
		ServLogin(w, r.WithContext(ctx))
	case !utils.ValidPassword(password):
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "invalid Password plese try again"})
		ServLogin(w, r.WithContext(ctx))
	case !exist:
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "this user oredy exist"})
		ServLogin(w, r.WithContext(ctx))
	default:
		userId, hash, err = db.GetUserHashByUsername(username)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return
		}
	}

	// if (!utils.ValidUsername(username) && !utils.ValidEmail(username)) || !utils.ValidPassword(password) || !exist {
	// 	ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "invalid credentials try again"})
	// 	ServLogin(w, r.WithContext(ctx))
	// 	return
	// }

	if !utils.CheckPassword(password, hash) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Wrong password try again"})
		ServLogin(w, r.WithContext(ctx))
		return
	}

	session := GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	err = db.UpdateUserSession(userId, session)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
