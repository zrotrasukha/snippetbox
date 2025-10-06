package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int
	Name            string
	Email           string
	Hashed_password []byte
	Created_at      time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) INSERT(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	stmt := `INSERT INTO users (name, email, password, created) 
					VALUES(?, ?, ?, UTC_TIMESTAMP());`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.user") {
				// custom error that we created in the errors.go, dayumn looking cool until now
				return ErrDupliacteEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := `SELECT id, password FROM users WHERE email = ?`
	err := m.DB.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
