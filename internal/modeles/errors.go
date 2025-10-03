package modeles

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDupliacteEmail = errors.New("models: duplicate email")
)


