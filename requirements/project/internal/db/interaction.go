package db

import (
	"database/sql"
	"errors"
	repo "forum/internal/repository"
)

func AlreadyExists(username, email string) (bool, error) {
	var count int
	err := repo.DB.QueryRow(repo.SELECT_USER_COUNT_BY_USERNAME_EMAIL, username, email).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}

func GetUserHashByUsername(username string) (int, string, error) {
	var hash string
	var id int
	err := repo.DB.QueryRow(repo.SELECT_USERID_PASSHASH_BY_USERNAME_EMAIL, username, username).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, hash, nil
		}
		return id, hash, err
	}
	return id, hash, nil
}

func GetUserInfo(user_id int) (repo.User, error) {
	var user repo.User

	err := repo.DB.QueryRow(repo.SELECT_USER_BY_ID, user_id).Scan(&user.Id, &user.Username,
		&user.Email, &user.Password_hash, &user.Created_at, &user.Updated_at)
	if err != nil {
		return user, err
	}
	return user, nil
}

func AddNewUser(username, email, hashedPass string) error {
	_, err := repo.DB.Exec(repo.INSERT_USERNAME_EMAIL_PASSHASH, username, email, hashedPass)
	return err
}

func GetUserHashById(id int) (string, error) {
	var hash string
	err := repo.DB.QueryRow(repo.SELECT_PASSHASH_BY_USERID, id).Scan(&hash)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

func DupplicatedUsername(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = ?`
	err := repo.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateUsernmae(id int, username string) error {
	query := `UPDATE users SET updated_at = DATETIME('now') , username = ? WHERE id = ?`
	res, err := repo.DB.Exec(query, username, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func DupplicatedEmail(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	err := repo.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateEmail(id int, email string) error {
	query := `UPDATE users SET updated_at = DATETIME('now') , email = ? WHERE id = ?`
	res, err := repo.DB.Exec(query, email, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err

}

func UpdatePassword(id int, password string) error {
	query := `UPDATE users SET updated_at = DATETIME('now'), password_hash = ? WHERE id = ?`
	res, err := repo.DB.Exec(query, password, id)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

func DeleteUser(user_id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := repo.DB.Exec(query, user_id)
	if err != nil {
		return err
	}
	return nil
}
