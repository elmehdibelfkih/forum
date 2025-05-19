package repository

const (
	// INSERT queries
	INSERT_NEW_POST         = `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	INSERT_NEW_COMMENT      = `INSERT INTO comments (user_id, post_id, comment) VALUES (?, ?, ?)`
	INSERT_NEW_LIKE_DISLIKE = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_dislike) VALUES (?, ?, ?, ?)`
	GET_POST_BYLIKES        = `
	SELECT 
    p.id,
    p.user_id,
    p.title,
    p.content,
    u.username AS publisher,
    IFNULL(GROUP_CONCAT(DISTINCT c.name), '') AS categories,
    COUNT(DISTINCT CASE WHEN ld.is_like = 1 THEN ld.user_id END) AS likes,
    COUNT(DISTINCT CASE WHEN ld.is_dislike = 1 THEN ld.user_id END) AS dislikes,
    IFNULL(GROUP_CONCAT(commenter.username || ':' || cm.comment ORDER BY datetime(cm.created_at) DESC), '') AS comments,
    p.created_at,
    p.updated_at
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN post_categories pc ON pc.post_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	LEFT JOIN likes_dislikes ld ON ld.post_id = p.id
	LEFT JOIN comments cm ON cm.post_id = p.id
	LEFT JOIN users commenter ON commenter.id = cm.user_id
	WHERE ld.user_id = ? AND ld.is_like = 1
	GROUP BY p.id
	ORDER BY p.created_at DESC;
	`
	GET_POST_BYOWNED = `
	SELECT 
    p.id,
    p.user_id,
    p.title,
    p.content,
    u.username AS publisher,
    IFNULL(GROUP_CONCAT(DISTINCT c.name), '') AS categories,
    COUNT(DISTINCT CASE WHEN ld.is_like = 1 THEN ld.user_id END) AS likes,
    COUNT(DISTINCT CASE WHEN ld.is_dislike = 1 THEN ld.user_id END) AS dislikes,
    IFNULL(GROUP_CONCAT(commenter.username || ':' || cm.comment ORDER BY datetime(cm.created_at) DESC), '') AS comments,
    p.created_at,
    p.updated_at
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN post_categories pc ON pc.post_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	LEFT JOIN likes_dislikes ld ON ld.post_id = p.id
	LEFT JOIN comments cm ON cm.post_id = p.id
	LEFT JOIN users commenter ON commenter.id = cm.user_id
	WHERE p.user_id = ?
	GROUP BY p.id
	ORDER BY p.created_at DESC;
	`
	// SELECT queries
	IS_POST_EXIST = `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	IS_LIKED      = `SELECT is_like FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	IS_DISLIKED   = `SELECT is_dislike FROM likes_dislikes WHERE user_id = ? AND post_id = ?`

	// UDDATE queries
	UPDATE_LIKE    = `UPDATE likes_dislikes SET is_like = ?, is_dislike = 0 WHERE user_id = ? AND post_id = ?`
	UPDATE_DISLIKE = `UPDATE likes_dislikes SET is_like = 0, is_dislike = ? WHERE user_id = ? AND post_id = ?`

	// todo: sort the comments by time
	SELECT_ALL_POSTS = `
  	SELECT
	p.id,
	p.user_id,
	p.title,
	p.content,
	u.username AS Publisher,
	IFNULL(GROUP_CONCAT(DISTINCT c.name), '') AS categories,
	COUNT(DISTINCT CASE WHEN ld.is_like = 1 THEN ld.user_id END) AS likes,
	COUNT(DISTINCT CASE WHEN ld.is_dislike = 1 THEN ld.user_id END) AS dislikes,
	IFNULL(GROUP_CONCAT(commenter.username || ':' || cm.comment ORDER BY datetime(cm.created_at) DESC), '') AS comments,
	p.created_at,
	p.updated_at
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN post_categories pc ON pc.post_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	LEFT JOIN likes_dislikes ld ON ld.post_id = p.id
	LEFT JOIN comments cm ON cm.post_id = p.id
	LEFT JOIN users commenter ON commenter.id = cm.user_id
	GROUP BY p.id
	ORDER BY p.created_at DESC;
`
)
