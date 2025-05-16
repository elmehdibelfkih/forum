package auth

import (
	"context"
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	utils "forum/internal/utils"
	"net/http"
)

func SwitchRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ServRegister(w, r)
	case http.MethodPost:
		SubmitRegister(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func ServRegister(w http.ResponseWriter, r *http.Request) {
	var errMap map[string]any
	if r.Context().Value(repo.ERROR_CASE) != nil {
		errMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "register.html", errMap)
}

func SubmitRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm_password := r.FormValue("confirm_password")

	if !utils.ValidUsername(username) || !utils.ValidEmail(email) || !utils.ValidPassword(password) || confirm_password != password { //TODO: it should be a better way
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "invalid credentials try again"})
		ServRegister(w, r.WithContext(ctx))
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	err = db.AddNewUser(username, email, hash)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" || err.Error() == "UNIQUE constraint failed: users.email" {
			ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "username or email alredy used"})
			ServRegister(w, r.WithContext(ctx))
			return
		} else {
			forumerror.InternalServerError(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
