package db

import (
	"database/sql"
	repo "forum/internal/repository"
	"strings"
)

func AddNewPost(user_id int, titel string, content string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_POST, user_id, titel, content)
	return err
}

// TODO: gel all catigories
// TODO: add offset and limit
func GetAllPostsInfo() (repo.PageData, error) {

	var data repo.PageData
	rows, err := repo.DB.Query(repo.SELECT_ALL_POSTS, 2)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var post repo.Post

		var categoriesStr, commentsStr string

		err := rows.Scan(
			&post.Id,
			&post.UserId,
			&post.Title,
			&post.Content,
			&post.Publisher,
			&categoriesStr,
			&post.Likes,
			&post.Dislikes,
			&commentsStr,
			&post.Created_at,
			&post.Updated_at,
		)
		if err != nil {
			return data, err
		}
		if categoriesStr != "" {
			post.Catigories = strings.Split(categoriesStr, ",")
		}
		if commentsStr != "" {
			comments := strings.Split(commentsStr, ",")
			for _, c := range comments {
				parts := strings.SplitN(c, ":", 2)
				if len(parts) == 2 {
					commentMap := map[string]string{
						parts[0]: parts[1],
					}
					post.Comments = append(post.Comments, commentMap)
				}
			}
		}
		post.IsEdited = post.Created_at != post.Updated_at

		data.Posts = append(data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

func IsPostExist(postId int) (bool, error) {
	var exists bool

	err := repo.DB.QueryRow(repo.IS_POST_EXIST, postId).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func AddNewComment(userId int, postId int, comment string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_COMMENT, userId, postId, comment)
	return err
}

func AddRemovePostLike(userId int, postId int) error {
	isLiked, err := IsPostLikedByUser(userId, postId)
	if err == sql.ErrNoRows {
		_, err = repo.DB.Exec(repo.INSERT_NEW_LIKE_DISLIKE, userId, postId, 1, 0)
		return err
	}
	if isLiked {
		res, err := repo.DB.Exec(repo.UPDATE_LIKE, 0, userId, postId)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err
	}
	res, err := repo.DB.Exec(repo.UPDATE_LIKE, 1, userId, postId)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func AddRemovePostDeslike(userId int, postId int) error {
	isDisike, err := IsPostDisikedByUser(userId, postId)
	if err == sql.ErrNoRows {
		_, err = repo.DB.Exec(repo.INSERT_NEW_LIKE_DISLIKE, userId, postId, 0, 1)
		return err
	}
	if isDisike {
		res, err := repo.DB.Exec(repo.UPDATE_DISLIKE, 0, userId, postId)
		if err != nil {
			return err
		}
		_, err = res.RowsAffected()
		return err
	}
	res, err := repo.DB.Exec(repo.UPDATE_DISLIKE, 1, userId, postId)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func IsPostLikedByUser(userId int, postId int) (bool, error) {
	var isLike bool

	err := repo.DB.QueryRow(repo.IS_LIKED, userId, postId).Scan(&isLike)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isLike == true {
		return true, nil
	}
	return false, nil
}

func IsPostDisikedByUser(userId int, postId int) (bool, error) {
	var isDisike bool

	err := repo.DB.QueryRow(repo.IS_DISLIKED, userId, postId).Scan(&isDisike)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isDisike == true {
		return true, nil
	}
	return false, nil
}
