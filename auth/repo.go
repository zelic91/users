package auth

import (
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

type User struct {
	ID             int64            `db:"id"`
	Username       string           `db:"username"`
	DisplayName    *string          `db:"display_name"`
	HashedPassword string           `db:"hashed_password"`
	App            string           `db:"app"`
	Email          *string          `db:"email"`
	CreatedAt      *strfmt.DateTime `db:"created_at"`
	UpdatedAt      *strfmt.DateTime `db:"updated_at"`
}

func NewRepo(db *sqlx.DB) Repo {
	return Repo{
		DB: db,
	}
}

func (r Repo) Create(user *User) (*User, error) {
	insertSQL := `
	INSERT INTO users (
		username,
		hashed_password,
		app
	) VALUES (
		$1,
		$2,
		$3
	)`

	_, err := r.DB.Exec(insertSQL, strings.ToLower(user.Username), user.HashedPassword, user.App)
	if err != nil {
		return nil, err
	}

	query := `
	SELECT * FROM users WHERE username = $1 LIMIT 1
	`
	result := User{}
	err = r.DB.Get(&result, query, user.Username)

	return &result, err
}

func (r Repo) GetByID(id int64) (*User, error) {
	query := "SELECT * FROM users WHERE id = ? LIMIT 1"
	result := User{}
	err := r.DB.Get(&result, query, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r Repo) GetByUsername(username string, app string) (*User, error) {
	query := "SELECT * FROM users WHERE username = $1 AND app = $2 LIMIT 1"
	result := User{}
	err := r.DB.Get(&result, query, strings.ToLower(username), strings.ToLower(app))
	if err != nil {
		return nil, err
	}

	return &result, nil
}
