package auth

import (
	forumerror "forum/internal/error"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		// No valid session, redirect to homepage
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionToken := sessionCookie.Value

	// Invalidate the session on server side
	hasSession, err := ResetUserSession(sessionToken)
	if err != nil {
		// Server-side error, do not proceed with cookie reset
		forumerror.TempErr(w, err, http.StatusInternalServerError)
		return
	}

	// Clear the session cookie on the client side
	clearCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Always in the past
		MaxAge:   -1,
	}
	http.SetCookie(w, clearCookie)

	// Handle invalid sessions
	if !hasSession {
		http.Error(w, "Session not found or already expired", http.StatusUnauthorized)
		return
	}

	// Redirect to home page after logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
