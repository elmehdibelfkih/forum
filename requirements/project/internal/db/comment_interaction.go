package db

import (
	"database/sql"
	repo "forum/internal/repository"
	"html"
	"html/template"
	"strings"
)

func AddNewComment(userId int, postId int, comment string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_COMMENT, userId, postId, comment)
	return err
}

func IsUserCanCommentToday(userId int) (bool, error) {
	var commentCount int
	err := repo.DB.QueryRow(repo.SELECT_TODAY_COMMENTS, userId).Scan(&commentCount)
	if err != nil {
		return false, err
	}
	return commentCount < repo.DAY_COMMENTS_LIMIT, nil
}

func GetCommentCount(postId int) (int, error) {
	var count int

	err := repo.DB.QueryRow(repo.GET_COMMENT_POST_COUNT, postId).Scan(&count)

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return count, nil
}

func GetCommentsByPostPaginated(postID, page, userID int) ([]repo.Comment, int, error) {
	offset := (page - 1) * repo.PAGE_COMMENT_QUANTITY

	rows, err := repo.DB.Query(repo.SELECT_COMMENT_BY_10, postID, repo.PAGE_COMMENT_QUANTITY, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []repo.Comment
	for rows.Next() {
		var username string
		var rawContent string
		var commentId int

		if err := rows.Scan(&username, &rawContent, &commentId); err != nil {
			return nil, 0, err
		}

		safe := html.EscapeString(rawContent)
		safe = strings.ReplaceAll(safe, "\r\n", "\n")
		safe = strings.ReplaceAll(safe, "\n", "<br>")

		c := repo.Comment{
			Username: username,
			Content:  template.HTML(safe),
		}
		if c.Username != "" {
			c.Initial = c.Username[:1]
		}
		c.CommentId = commentId
		c.PostId = postID
		err = MakeCommentMetadata(&c, userID)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, c)
	}
	var total int

	err = repo.DB.QueryRow(repo.GET_COMMENT_POST_COUNT, postID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// TODO:
func MakeCommentMetadata(comment *repo.Comment, userId int) error {
	var err error

	comment.IsCommentLikedByUser, err = IsCommentLikedByUser(userId, comment.CommentId)
	if err == sql.ErrNoRows {
		err = nil
	}
	comment.IsCommentDislikedByUser, err = IsCommentDisikedByUser(userId, comment.CommentId)
	if err == sql.ErrNoRows {
		err = nil
	}

	// println(comment.PostId)
	// println(comment.CommentId)

	comment.CommentLikes = 10
	comment.CommentDislikes = 10
	// err = 10
	// err = 10

	return err
}

// TODO:
func IsCommentExist(postId, commentId int) (bool, error) {
	var exists bool

	err := repo.DB.QueryRow(repo.IS_COMMENT_EXIST, postId, commentId).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// TODO:
func AddRemoveCommentLike(userId, commentId int) error {
	isLiked, err := IsCommentLikedByUser(userId, commentId)
	if err == sql.ErrNoRows {
		_, err = repo.DB.Exec(repo.INSERT_COMMENT_LIKE_DISLIKE, userId, commentId, 1, 0)
		return err
	}
	if isLiked {
		res, err := repo.DB.Exec(repo.UPDATE_COMMENT_LIKE, 0, userId, commentId)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err
	}
	res, err := repo.DB.Exec(repo.UPDATE_COMMENT_LIKE, 1, userId, commentId)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

// TODO:
func AddRemoveCommentDislike(userId, commentId int) error {
	isDisike, err := IsCommentDisikedByUser(userId, commentId)
	if err == sql.ErrNoRows {
		_, err = repo.DB.Exec(repo.INSERT_COMMENT_LIKE_DISLIKE, userId, commentId, 0, 1)
		return err
	}
	if isDisike {
		res, err := repo.DB.Exec(repo.UPDATE_COMMENT_DISLIKE, 0, userId, commentId)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err
	}
	res, err := repo.DB.Exec(repo.UPDATE_COMMENT_DISLIKE, 1, userId, commentId)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

// TODO:
func IsCommentLikedByUser(userId, commentId int) (bool, error) {
	var isLike bool

	err := repo.DB.QueryRow(repo.IS_COMMENT_LIKED, userId, commentId).Scan(&isLike)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isLike {
		return true, nil
	}
	return false, nil
}

// TODO:
func IsCommentDisikedByUser(userId, commentId int) (bool, error) {
	var isDisike bool

	err := repo.DB.QueryRow(repo.IS_COMMENT_DISLIKED, userId, commentId).Scan(&isDisike)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isDisike {
		return true, nil
	}
	return false, nil
}
