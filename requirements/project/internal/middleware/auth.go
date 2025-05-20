package middleware

import (
	"context"
	"fmt"
	auth "forum/internal/auth"
	db "forum/internal/db"
	forumerror "forum/internal/error"
	repo "forum/internal/repository"
	"net/http"
)

func AuthMidleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		fmt.Printf("%s",r.URL.Path)
		if r.URL.Path == "/filterby" { // CASE WHERE THE USER USE FILERT I WILL CHECK IF HE IS LOGING OR NOT TO CONTROL THE PAGE WEB DETAIL !
			if err == nil && sessionCookie.Value != "" {
				// Try to get session from DATABASE !!
				userID, exist, err := db.SelectUserSession(sessionCookie.Value)
				if err != nil {
					forumerror.InternalServerError(w, r, err)
					return
				}
				if exist {
					// Inject user ID and auth WILL Be true he is login !!
					ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, userID)
					ctx = context.WithValue(ctx, repo.PUBLIC, true)
					next(w, r.WithContext(ctx))
					return
				}
			}
			// If no cookie or session not found: still allow, but mark as not authenticated "NOT LOGIN"
			ctx := context.WithValue(r.Context(), repo.PUBLIC, false)
			next(w, r.WithContext(ctx))
			return
		}
		 
		// At mdileware i will check if he is login 
		if err != nil || sessionCookie.Value == "" {
			auth.ServLogin(w, r)
			return
		}

		user_id, exist, err := db.SelectUserSession(sessionCookie.Value)
		
		if err != nil {
		forumerror.InternalServerError(w,r, err)
		return
		}

		if !exist {
			auth.ServLogin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), repo.USER_ID_KEY, user_id) //avoid collisions
		ctx = context.WithValue(ctx, repo.PUBLIC, true) // HER IS TRACK THE USER IF HE IS LOGIN OR NOT !
		next(w, r.WithContext(ctx))
	}
}
