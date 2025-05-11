package repository

const (
	INSERT_NEW_POST = `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
)
