package db

import (
	"fmt"
	repo "forum/internal/repository"
	"net/http"
	"strings"
)


func Getposbytlikes(r *http.Request) (repo.PageData, error) {
	user_id := r.Context().Value(repo.USER_ID_KEY).(int)
	fmt.Printf("%d", user_id)
	rows, err := repo.DB.Query(repo.GET_POST_BYLIKES, user_id)
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

		data.Posts = append(data.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return data, err
	}

	return data, nil
}
