package repository

const (
	// INSERT queries
	INSERT_NEW_POST         = `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	INSERT_NEW_COMMENT      = `INSERT INTO comments (user_id, post_id, comment) VALUES (?, ?, ?)`
	INSERT_NEW_LIKE_DISLIKE = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_dislike) VALUES (?, ?, ?, ?)`

	//SELECT queries
	SELECT_ALL_POSTS      = `SELECT * FROM posts ORDER BY datetime(created_at) DESC`
	SELECT_LIKES_COUNT    = `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND is_like = 1`
	SELECT_DISLIKES_COUNT = `SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND is_dislike = 1`
	SELECT_COMMENTS       = `SELECT users.username, comments.comment FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY datetime(comments.created_at) DESC`
	IS_POST_EXIST         = `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	IS_LIKED              = `SELECT is_like FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	IS_DISLIKED           = `SELECT is_dislike FROM likes_dislikes WHERE user_id = ? AND post_id = ?`

	// UDDATE queries
	UPDATE_LIKE    = `UPDATE likes_dislikes SET is_like = ?, is_dislike = 0 WHERE user_id = ? AND post_id = ?`
	UPDATE_DISLIKE = `UPDATE likes_dislikes SET is_like = 0, is_dislike = ? WHERE user_id = ? AND post_id = ?`
)

