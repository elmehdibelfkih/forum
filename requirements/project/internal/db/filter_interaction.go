package db

import (
	"database/sql"
	repo "forum/internal/repository"
	"strings"
)

func Getposbytlikes(userId int, page int) (repo.PageData, error) {
	var data repo.PageData

	rows, err := repo.DB.Query(repo.GET_POST_BYLIKES, userId, repo.PAGE_POSTS_QUANTITY, (page-1)*repo.PAGE_POSTS_QUANTITY)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var post repo.Post
		var categoriesStr, commentsStr string

		err := rows.Scan(
			&post.Id,
			&post.Publisher,
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
			comments := strings.SplitSeq(commentsStr, ",")
			for c := range comments {
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
		post.IsLikedByUser = true
		data.Posts = append(data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}

func Getpostbyowner(userId int, page int) (repo.PageData, error) {
	var data repo.PageData

	rows, err := repo.DB.Query(repo.GET_POST_BYOWNED, userId, repo.PAGE_POSTS_QUANTITY, (page-1)*repo.PAGE_POSTS_QUANTITY)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var post repo.Post
		var categoriesStr, commentsStr string

		err := rows.Scan(
			&post.Id,
			&post.Publisher,
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
			comments := strings.SplitSeq(commentsStr, ",")
			for c := range comments {
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
		post.IsLikedByUser, err = IsPostLikedByUser(userId, post.Id)
		if err != nil && err != sql.ErrNoRows {
			return data, err
		}
		post.IsDislikedByUser, err = IsPostDisikedByUser(userId, post.Id)
		if err != nil && err != sql.ErrNoRows {
			return data, err
		}
		data.Posts = append(data.Posts, post)
	}
	if err := rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}

func GePostbycategory(category string, page int, userId int) (repo.PageData, error) {
	var data repo.PageData

	rows, err := repo.DB.Query(repo.GET_POST_BYCATEGORY, category, repo.PAGE_POSTS_QUANTITY, (page-1)*repo.PAGE_POSTS_QUANTITY)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var post repo.Post
		var categoriesStr, commentsStr string

		err := rows.Scan(
			&post.Id,
			&post.Publisher,
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
			comments := strings.SplitSeq(commentsStr, ",")
			for c := range comments {
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
		post.IsLikedByUser, err = IsPostLikedByUser(userId, post.Id)
		if err != nil && err != sql.ErrNoRows {
			return data, err
		}
		post.IsDislikedByUser, err = IsPostDisikedByUser(userId, post.Id)
		if err != nil && err != sql.ErrNoRows {
			return data, err
		}
		data.Posts = append(data.Posts, post)
	}
	if err := rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}

