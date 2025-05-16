package db

import (
	"database/sql"
	repo "forum/internal/repository"
)

func AddNewPost(user_id int, titel string, content string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_POST, user_id, titel, content)
	return err
}

// TODO: gel all catigories
func GetAllPostsInfo() (repo.PageData, error) {
	var data repo.PageData
	var post repo.Post

	rows, err := repo.DB.Query(repo.SELECT_ALL_POSTS)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Created_at, &post.Updated_at)
		if err != nil {
			return data, err
		}
		post.IsEdited = post.Updated_at != post.Created_at
		// TODO: fix time format
		post.Publisher, err = GetUserNameById(post.UserId)
		if err != nil {
			return data, err
		}
		post.Likes, err = GetPostLikes(post.Id)
		if err != nil {
			return data, err
		}
		post.Dislikes, err = GetPostDislikes(post.Id)
		if err != nil {
			return data, err
		}
		post.Comments, err = GetPostComments(post.Id)
		if err != nil {
			return data, err
		}
		data.Posts = append(data.Posts, post)
	}
	if err = rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}

func GetPostLikes(postId int) (int, error) {
	var likes int

	err := repo.DB.QueryRow(repo.SELECT_LIKES_COUNT, postId).Scan(&likes)
	if err != nil {
		return likes, err
	}
	return likes, nil
}

func GetPostDislikes(postId int) (int, error) {
	var dislikes int

	err := repo.DB.QueryRow(repo.SELECT_DISLIKES_COUNT, postId).Scan(&dislikes)
	if err != nil {
		return dislikes, err
	}
	return dislikes, nil
}

func GetPostComments(postId int) ([]map[string]string, error) {
	var comments []map[string]string
	var userNameTmp, commentTmp string

	rows, err := repo.DB.Query(repo.SELECT_COMMENTS, postId)
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userNameTmp, &commentTmp)
		if err != nil {
			return comments, err
		}
		toAppend := map[string]string{
			userNameTmp: commentTmp,
		}
		comments = append(comments, toAppend)
	}

	if err = rows.Err(); err != nil {
		return comments, err
	}

	return comments, nil
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
	var isLike int

	err := repo.DB.QueryRow(repo.IS_LIKED, userId, postId).Scan(&isLike)
	println("isLike : ", isLike)
	println("hani",postId, userId)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isLike == 1 {
		return true, nil
	}
	return false, nil
}

func IsPostDisikedByUser(userId int, postId int) (bool, error) {
	var isDisike int

	err := repo.DB.QueryRow(repo.IS_DISLIKED, userId, postId).Scan(&isDisike)
	if err == sql.ErrNoRows {
		return false, err
	} else if err != nil {
		return false, err
	} else if isDisike == 1 {
		return true, nil
	}
	return false, nil
}
