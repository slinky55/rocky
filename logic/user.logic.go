package logic

import (
	"github.com/slinky55/rocky/db"
)

func UserExists(userId string) bool {
	rows, err := db.Conn.Queryx("SELECT user_id FROM users WHERE user_id = ?", userId)
	if err != nil {
		return true
	}

	if rows.Next() {
		return true
	}

	return false
}

func CreateUser(userId string) error {
	_, err := db.Conn.Exec("INSERT INTO users (user_id, created_at) VALUES (?, unixepoch())", userId)
	return err
}

func VerifyUser(userId string) error {
	_, err := db.Conn.Exec("UPDATE users SET verified = ? WHERE user_id = ?", true, userId)
	return err
}

func SetUserPublicKey(userId string, publicKey string) error {
	_, err := db.Conn.Exec("UPDATE users SET public_key = ? WHERE user_id = ?", publicKey, userId)
	return err
}

func GetUserPublicKey(userId string) (string, error) {
	rows, err := db.Conn.Query("SELECT public_key FROM users WHERE user_id = ?", userId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var publicKey string
		if err := rows.Scan(&publicKey); err != nil {
			return "", err
		}
		return publicKey, nil
	} else {
		return "", nil
	}
}
