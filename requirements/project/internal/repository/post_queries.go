package repository

const (
	INSERT_NEW_POST         = `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	INSERT_NEW_COMMENT      = `INSERT INTO comments (user_id, post_id, comment) VALUES (?, ?, ?)`
	INSERT_NEW_LIKE_DESLIKE = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_deslike) VALUES (?, ?, ?, ?)`
	// INSERT_NEW_DESLIKE = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_deslike) VALUES (?, ?, ?, ?)`

	SELECT_ALL_POSTS      = `SELECT * FROM posts`
	SELECT_LIKES_COUNT    = `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND is_like = 1`
	SELECT_DESLIKES_COUNT = `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND likes_dislikes = 1`
	SELECT_COMMENTS       = `SELECT users.username, comments.comment FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = ?`
	IS_POST_EXIST         = `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	IS_LIKED              = `SELECT is_like FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	IS_DESLIKED           = `SELECT is_deslike FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
)
