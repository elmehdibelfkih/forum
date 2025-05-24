package db

import (
	"database/sql"
	"strings"

	repo "forum/internal/repository"
)

func Getposbytlikes(userId int) (repo.PageData, error) {
	rows, err := repo.DB.Query(repo.GET_POST_BYLIKES, userId)
	if err != nil {
		return repo.PageData{}, err
	}
	defer rows.Close()

	var data repo.PageData

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

		// Split categories if they exist
		if categoriesStr != "" {
			post.Catigories = strings.Split(categoriesStr, ",")
		}

		// Split and parse comments if they exist
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

		// Check if the post was edited
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsLikedByUser, err = IsPostLikedByUser(post.UserId, post.Id)
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
func Getpostbyowner(userId int) (repo.PageData, error) {

	rows, err := repo.DB.Query(repo.GET_POST_BYOWNED, userId)
	if err != nil {
		return repo.PageData{}, err
	}
	defer rows.Close()

	var data repo.PageData

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

		// Split categories if they exist
		if categoriesStr != "" {
			post.Catigories = strings.Split(categoriesStr, ",")
		}

		// Split and parse comments if they exist
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

		// Check if the post was edited
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsLikedByUser, err = IsPostLikedByUser(post.UserId, post.Id)
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

func GePostbycategory(category string) (repo.PageData, error) {

	// i want to get all post by the category there the user want
	// i have post_category and i have a table of category
	// and when user creat post i register it at category_post...!!
	rows, err := repo.DB.Query(repo.GET_POST_BYCATEGORY, category)
	if err != nil {
		return repo.PageData{}, err
	}
	defer rows.Close()

	var data repo.PageData

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

		// Split categories if they exist
		if categoriesStr != "" {
			post.Catigories = strings.Split(categoriesStr, ",")
		}

		// Split and parse comments if they exist
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

		// Check if the post was edited
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsEdited = post.Created_at != post.Updated_at
		post.IsLikedByUser, err = IsPostLikedByUser(post.UserId, post.Id)
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
