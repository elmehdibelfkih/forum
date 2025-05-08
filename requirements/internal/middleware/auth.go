package middleware

import (
	"context"
	auth "forum/internal/auth"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
)

func AuthMidleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil || sessionCookie.Value == "" {
			auth.ServLogin(w, r)
			return
		}

		user_id, exist, err := auth.SelectUserSession(sessionCookie.Value)

		if err != nil {
			forumerror.TempErr(w, err, http.StatusInternalServerError)
		}

		if !exist {
			auth.ServLogin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, user_id) //avoid collisions
		next(w, r.WithContext(ctx))
	}
}
