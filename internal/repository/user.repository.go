package repository

import (
	"database/sql"

	"github.com/NekKkMirror/erm-tz/internal/model"
)

// UserRepository provides methods for managing user records in the database.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates and returns a new UserRepository instance tied to the given database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save inserts a new user record into the "users" table with the provided user data.
// It takes a User object as input and writes its fields (nickname, email, verification status, created time)
// to the corresponding columns in the database.
func (ur *UserRepository) Save(u *model.User) error {
	q := "INSERT INTO users (nickname, email, verified, created_at) VALUES ($1, $2, $3, $4)"
	_, err := ur.db.Exec(q, u.Nickname, u.Email, u.Verified, u.CreatedAt)
	return err
}

// VerifyEmail updates the "verified" status of a user in the "users" table to "true"
// based on the provided email address. It marks the user's email as verified in the database.
func (ur *UserRepository) VerifyEmail(email string) error {
	q := "UPDATE users SET verified = true WHERE email = $1"
	_, err := ur.db.Exec(q, email)
	return err
}
