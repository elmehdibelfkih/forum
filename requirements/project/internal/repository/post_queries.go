package repository

const (
	// INSERT queries
	INSERT_NEW_POST = `
                            INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?);
                            UPDATE post_metadata SET post_count = post_count + 1;`
	INSERT_NEW_COMMENT         = `INSERT INTO comments (user_id, post_id, comment) VALUES (?, ?, ?)`
	INIT_POST_META_DATA        = `INSERT OR IGNORE INTO post_metadata (id, post_count) VALUES (1, 0);`
	INIT_FIELDS_QUERY          = `INSERT OR IGNORE INTO categories (name) VALUES (?)`
	INIT_POST_CATEGORIES_COUNT = `INSERT OR IGNORE INTO post_categories_count (category_id, post_count) VALUES (?, 0)`
	INSERT_NEW_LIKE_DISLIKE    = `INSERT INTO likes_dislikes (user_id, post_id, is_like, is_dislike) VALUES (?, ?, ?, ?)`
  // todo: update the post categories count on delete
	MAP_POSTS_WITH_CATEGORY    = `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?);
                                UPDATE post_categories_count SET post_count = post_count + 1 WHERE category_id = ?;`

	// DELETE queries
	DELETE_POST = `
                DELETE FROM posts WHERE id = ?;
                UPDATE post_metadata SET post_count = post_count - 1;`

	// SELECT queries
	IS_POST_EXIST      = `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	IS_LIKED           = `SELECT is_like FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	IS_DISLIKED        = `SELECT is_dislike FROM likes_dislikes WHERE user_id = ? AND post_id = ?`
	SELECT_TODAY_POSTS = `SELECT COUNT(*) FROM posts WHERE user_id = ?  AND created_at >= DATE('now')` // FIXME: fix the time
	SELECT_CATEGORY_ID = `SELECT id FROM categories WHERE name = ?`

	// UDDATE queries
	UPDATE_LIKE    = `UPDATE likes_dislikes SET is_like = ?, is_dislike = 0 WHERE user_id = ? AND post_id = ?`
	UPDATE_DISLIKE = `UPDATE likes_dislikes SET is_like = 0, is_dislike = ? WHERE user_id = ? AND post_id = ?`
	GET_POST_COUNT = `SELECT post_count FROM post_metadata`

  GET_POST_COUNT_BY_CAT = `SELECT pcc.post_count FROM post_categories_count pcc JOIN categories c ON pcc.category_id = c.id WHERE c.name = ?`
  

	// JOIN queries
	GET_POST_BYLIKES = `
	 SELECT 
    p.id,
    p.user_id,
    p.title,
    p.content,
    u.username AS publisher,
    IFNULL(pc.categories, '') AS categories,
    IFNULL(ld.likes, 0) AS likes,
    IFNULL(ld.dislikes, 0) AS dislikes,
    IFNULL(cm.comments, '') AS comments,
    p.created_at,
    p.updated_at
  FROM posts p
  JOIN users u ON u.id = p.user_id
  LEFT JOIN (
  SELECT 
    pc.post_id,
  GROUP_CONCAT(DISTINCT c.name) AS categories
  FROM post_categories pc
  JOIN categories c ON c.id = pc.category_id
    GROUP BY pc.post_id
    ) pc ON pc.post_id = p.id
  LEFT JOIN (
  SELECT 
    ld.post_id,
    COUNT(DISTINCT CASE WHEN ld.is_like = 1 THEN ld.user_id END) AS likes,
    COUNT(DISTINCT CASE WHEN ld.is_dislike = 1 THEN ld.user_id END) AS dislikes
  FROM likes_dislikes ld
    GROUP BY ld.post_id 
    ) ld ON ld.post_id = p.id
  LEFT JOIN (
  SELECT 
    cm.post_id,
    GROUP_CONCAT(DISTINCT commenter.username || ':' || cm.comment ORDER BY datetime(cm.created_at) DESC) AS comments
  FROM comments cm
  JOIN users commenter ON commenter.id = cm.user_id
    GROUP BY cm.post_id
    ) cm ON cm.post_id = p.id
  WHERE p.id IN (
      SELECT post_id FROM likes_dislikes WHERE user_id = ? AND is_like = 1
    )
  ORDER BY p.created_at DESC
  LIMIT ? OFFSET ?;
`
	GET_POST_BYOWNED = `
    SELECT 
    p.id,
    p.user_id,
    p.title,
    p.content,
    u.username AS publisher,
    IFNULL(pc.categories, '') AS categories,
    IFNULL(ld.likes, 0) AS likes,
    IFNULL(ld.dislikes, 0) AS dislikes,
    IFNULL(cm.comments, '') AS comments,
    p.created_at,
    p.updated_at
    FROM posts p
    JOIN users u ON u.id = p.user_id
    LEFT JOIN (
    SELECT 
        pc.post_id,
        GROUP_CONCAT(DISTINCT c.name) AS categories
    FROM post_categories pc
    JOIN categories c ON c.id = pc.category_id
    GROUP BY pc.post_id
    ) pc ON pc.post_id = p.id
    LEFT JOIN (
    SELECT 
        ld.post_id,
        COUNT(DISTINCT CASE WHEN ld.is_like = 1 THEN ld.user_id END) AS likes,
        COUNT(DISTINCT CASE WHEN ld.is_dislike = 1 THEN ld.user_id END) AS dislikes
    FROM likes_dislikes ld
    GROUP BY ld.post_id
    ) ld ON ld.post_id = p.id
    LEFT JOIN (
    SELECT 
        cm.post_id,
        GROUP_CONCAT(commenter.username || ':' || cm.comment ORDER BY datetime(cm.created_at) DESC) AS comments
    FROM comments cm
    JOIN users commenter ON commenter.id = cm.user_id
    GROUP BY cm.post_id
    ) cm ON cm.post_id = p.id
    WHERE p.user_id = ?
      ORDER BY p.created_at DESC
        LIMIT ? OFFSET ?;
	`
	GET_POST_BYCATEGORY = `
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
    WHERE c.name = ?
    GROUP BY p.id
    ORDER BY p.created_at DESC
    LIMIT ? OFFSET ?;
  `

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
         GROUP_CONCAT(users.username || ':' || comments.comment
                      ORDER BY datetime(comments.created_at) DESC) AS comments
  FROM comments
  JOIN users ON users.id = comments.user_id
  GROUP BY comments.post_id
  ) AS com ON com.post_id = posts.id

  ORDER BY posts.created_at DESC
  LIMIT ? OFFSET ?;
  `
)
