package handler

import (
	"context"
	auth "forum/internal/auth"
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	utils "forum/internal/utils"
	"net/http"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	user, err := db.GetUserInfo(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "profile.html", user) // when u excute 2 template the get concatinated one in top of the other
}

func UpddateProfile(w http.ResponseWriter, r *http.Request) {
	var confMap = make(map[string]any)
	value := r.PathValue("value")
	if r.Context().Value(repo.ERROR_CASE) != nil {
		confMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}

	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	user, err := db.GetUserInfo(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	confMap["Username"] = user.Username

	switch value {
	case "username":
		confMap["username"] = true
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", confMap)
		return
	case "email":
		confMap["email"] = true
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", confMap)
		return
	case "password":
		confMap["password"] = true
		repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "update.html", confMap)
		return
	default:
		forumerror.BadRequest(w, r)
	}
}

func SaveChanges(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("value") {
	case "username":
		SaveUsername(w, r)
		return
	case "email":
		SaveEmail(w, r)
		return
	case "password":
		SavePassword(w, r)
		return
	default:
		http.Error(w, "bad req", http.StatusBadRequest)
		return
	}
}

func SaveUsername(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	new_username := r.FormValue("username")
	password := r.FormValue("current")
	if !utils.ValidUsername(new_username) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Please enter a valid username"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}

	hash, err := db.GetUserHashById(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if !utils.CheckPassword(password, hash) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Wrong password"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}

	dupp, err := db.DupplicatedUsername(new_username)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if dupp {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Username Alredy exists try again"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}

	err = db.UpdateUsernmae(user_id, new_username)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func SaveEmail(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	new_email := r.FormValue("email")
	password := r.FormValue("current")
	if !utils.ValidEmail(new_email) {
		ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, map[string]any{"Error": true, "Message": "Invalid email try again"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}
	hash, err := db.GetUserHashById(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if !utils.CheckPassword(password, hash) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Wrong password"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}

	dupp, err := db.DupplicatedEmail(new_email)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if dupp {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Email Alredy exists try again"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}

	err = db.UpdateEmail(user_id, new_email)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func SavePassword(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	current := r.FormValue("current")
	new := r.FormValue("new")
	confirm := r.FormValue("confirm")
	hash, err := db.GetUserHashById(user_id)

	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if current == new || !utils.ValidPassword(new) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "You used an Old password"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}
	if !utils.CheckPassword(current, hash) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Wrong password"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}
	if new != confirm {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Please Confirm Your password"})
		UpddateProfile(w, r.WithContext(ctx))
		return
	}
	new_hash, err := utils.HashPassword(new)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	err = db.UpdatePassword(user_id, new_hash)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func ServeDelete(w http.ResponseWriter, r *http.Request) {
	var confMap = make(map[string]any)

	if r.Context().Value(repo.ERROR_CASE) != nil {
		confMap = r.Context().Value(repo.ERROR_CASE).(map[string]any)
	}
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	user, err := db.GetUserInfo(user_id)
	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	confMap["Username"] = user.Username
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "delete.html", confMap)
}

func DeleteConfirmation(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	password := r.FormValue("password")
	hash, err := db.GetUserHashById(user_id)

	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
	if !utils.CheckPassword(password, hash) {
		ctx := context.WithValue(r.Context(), repo.ERROR_CASE, map[string]any{"Error": true, "Message": "Wrong password"})
		ServeDelete(w, r.WithContext(ctx))
		return
	}
	auth.LogoutHandler(w, r)
	err = db.DeleteUser(user_id)

	if err != nil {
		forumerror.InternalServerError(w, r, err)
		return
	}
}
