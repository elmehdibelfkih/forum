package db

import repo "forum/internal/repository"

func AddNewPost(user_id int, titel string, content string) error {
	_, err := repo.DB.Exec(repo.INSERT_NEW_POST, user_id, titel, content)
	return err
}
