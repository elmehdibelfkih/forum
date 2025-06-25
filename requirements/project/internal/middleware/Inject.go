package middleware

import (
	"context"
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
	"time"
)

func InjectUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")

		if err != nil || sessionCookie.Value == "" {
			ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, -1)
			next(w, r.WithContext(ctx))
			return
		}
		userId, expires_at, exist, err := db.SelectUserSession(sessionCookie.Value)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return
		}

		// avoid cookie extends from the browser utc for the sql
		if expires_at.Sub(time.Now().UTC()) <= 0 {
			ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, -1)
			next(w, r.WithContext(ctx))
			return
		}

		ctx := r.Context()
		if !exist {
			ctx = context.WithValue(ctx, repo.USER_ID_KEY, -1)
			next(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(ctx, repo.USER_ID_KEY, userId)

		usrName, err := db.GetUserNameById(userId)
		if err != nil {
			forumerror.InternalServerError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, repo.USER_NAME, usrName)
		next(w, r.WithContext(ctx))

	}
}
