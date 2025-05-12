package db

import (
	"database/sql"
	repo "forum/internal/repository"
)

func AddNewPost(user_id int, titel string, content string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_POST, user_id, titel, content)
	return err
}

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
		post.Publisher, err = GetUserNameById(post.UserId)
		if err != nil {
			return data, err
		}
		post.Likes, err = GetPostLikes(post.Id)
		if err != nil {
			return data, err
		}
		post.Deslikes, err = GetPostDeslikes(post.Id)
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

func GetPostDeslikes(postId int) (int, error) {
	var deslikes int

	err := repo.DB.QueryRow(repo.SELECT_DESLIKES_COUNT, postId).Scan(&deslikes)
	if err != nil {
		return deslikes, err
	}
	return deslikes, nil
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
			"username": userNameTmp,
			"comment":  commentTmp,
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
