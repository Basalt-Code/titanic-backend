package repositories

import "fmt"

var (
	columnMapping = map[string]string{
		"users_nickname_key": "nickname",
		"users_email_key":    "email",
	}
)

type ErrDuplicateField struct {
	Field string
}

func (err ErrDuplicateField) Error() string {
	return fmt.Sprintf("%s уже занят", columnMapping[err.Field])
}
