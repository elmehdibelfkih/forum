package repository

const (
	// insert queries
	INSERT_NEW_SESSION             = `INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, DATETIME('now', '+1 hour'))`
	INSERT_NEW_USER                = `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
	INSERT_USERNAME_EMAIL_PASSHASH = `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`

	// select queries
	SELECT_USER_BY_ID                        = `SELECT * FROM users WHERE id = ?`
	SELECT_USER_BY_SESSION_TOKEN             = `SELECT user_id FROM sessions WHERE session_token = ?`
	SELECT_USER_COUNT_BY_USERNAME_EMAIL      = `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	SELECT_USERID_PASSHASH_BY_USERNAME_EMAIL = `SELECT id,password_hash FROM users WHERE username = ? OR email = ?`
	SELECT_PASSHASH_BY_USERID                = `SELECT password_hash FROM users WHERE id = ?`

	// update queries
	UPDATE_SESSION_EXPIRING_TIME = `UPDATE sessions SET expires_at = DATETIME('now', '+1 hour'), session_token = ? WHERE user_id = ?`
	RESET_USER_SESSION_TOKEN     = `UPDATE sessions SET session_token = NULL WHERE session_token = ?`
)
