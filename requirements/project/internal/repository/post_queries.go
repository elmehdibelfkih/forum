package repository

const (
	// INSERT queries
	INSERT_NEW_POST         = `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
	INSERT_NEW_COMMENT      = `INSERT INTO comments (user_id, post_id, comment) VALUES (?, ?, ?)`
	INSERT_NEW_LIKE_DISLIKE = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_dislike) VALUES (?, ?, ?, ?)`
	MAP_POSTS_WITH_CATEGORY = `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`

	//SELECT queries
	IS_POST_EXIST      = `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	IS_LIKED           = `SELECT is_like FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	IS_DISLIKED        = `SELECT is_dislike FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	SELECT_CATEGORY_ID = `SELECT id FROM categories WHERE name = ?`

	// UDDATE queries
	UPDATE_LIKE    = `UPDATE likes_dislikes SET is_like = ?, is_dislike = 0 WHERE user_id = ? AND post_id = ?`
	UPDATE_DISLIKE = `UPDATE likes_dislikes SET is_like = 0, is_dislike = ? WHERE user_id = ? AND post_id = ?`

	// todo: sort the comments by time
	SELECT_ALL_POSTS = `
SELECT
  posts.id,
  posts.user_id,
  posts.title,
  posts.content,
  users.username,
  IFNULL(cat.categories, '') AS categories,
  COALESCE(likes.likes, 0) AS likes,
  COALESCE(dislikes.dislikes, 0) AS dislikes,
  IFNULL(com.comments, '') AS comments,
  posts.created_at,
  posts.updated_at
FROM posts
JOIN users ON users.id = posts.user_id

LEFT JOIN (
  SELECT post_categories.post_id,
         GROUP_CONCAT(DISTINCT categories.name) AS categories
  FROM post_categories
  JOIN categories ON categories.id = post_categories.category_id
  GROUP BY post_categories.post_id
) AS cat ON cat.post_id = posts.id

LEFT JOIN (
  SELECT post_id,
         COUNT(DISTINCT user_id) AS likes
  FROM likes_dislikes
  WHERE is_like = 1
  GROUP BY post_id
) AS likes ON likes.post_id = posts.id

LEFT JOIN (
  SELECT post_id,
         COUNT(DISTINCT user_id) AS dislikes
  FROM likes_dislikes
  WHERE is_dislike = 1
  GROUP BY post_id
) AS dislikes ON dislikes.post_id = posts.id

LEFT JOIN (
  SELECT comments.post_id,
         GROUP_CONCAT(DISTINCT users.username || ':' || comments.comment
                      ORDER BY datetime(comments.created_at) DESC) AS comments
  FROM comments
  JOIN users ON users.id = comments.user_id
  GROUP BY comments.post_id
) AS com ON com.post_id = posts.id

ORDER BY posts.created_at DESC;

`
)
