package modeles

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int
	Name       string
	Email      string
	Hashed_password []byte 
	Created_at time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) INSERT (name , email, password string) error {
	return nil 
}

func (u *UserModel) Authenticate (email, password string) (int, error) {
	return 0, nil
}
func (u *UserModel) Exists (id int) (bool, error) {
	return false, nil
}


