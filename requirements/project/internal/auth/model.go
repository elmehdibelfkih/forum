package auth

import (
	"database/sql"
	repo "forum/internal/repository"
	"log"

	"errors"
)

func SelectUserSession(session_id string) (int, bool, error) {
	var user_id int

	err := repo.DB.QueryRow(repo.SELECT_USER_BY_SESSION_TOKEN, session_id).Scan(&user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_id, false, nil // not logged in
		}
		return user_id, false, err // err in database
	}
	return user_id, true, nil // logged in
}

func UpdateUserSession(id int, token string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return err
	}

	res, err := tx.Exec(repo.UPDATE_SESSION_EXPIRING_TIME, token, id)
	if err != nil {
		log.Println("Update error:", err)
		tx.Rollback()
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows:", err)
		tx.Rollback()
		return err
	}

	if count > 0 {
		return tx.Commit() // if update succed
	}

	_, err = tx.Exec(repo.INSERT_NEW_SESSION, id, token)
	if err != nil {
		log.Println("Insert error:", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ResetUserSession(session_id string) (bool, error) {
	_, err := repo.DB.Exec(repo.RESET_USER_SESSION_TOKEN, session_id)
	if err != nil {
		return false, err
	}
	return true, nil
}
