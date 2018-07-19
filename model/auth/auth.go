package auth

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// Credential store user credential
type Credential struct {
	UserAlias string `json:"userAlias"`
	Password  string `json:"password"`
}

// DB is *sql.DB store auth information
type DB struct {
	*sql.DB
}

// Insert insert credential for authentication
func (db *DB) Insert(userAlias, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO AUTHENTICATION (
      useralias,
      password
    ) values (?, ?)`,
		userAlias, hashedPassword,
	)
	return err
}

// Authenticate will lookup userAlias in sqlite3 database and
// compare the hashed password store in database
func (db *DB) Authenticate(userAlias, password string) error {
	type auth struct {
		userAlias string
		password  []byte
	}
	a := new(auth)

	if err := db.QueryRow(
		`SELECT * FROM AUTHENTICATION WHERE useralias=?`,
		userAlias,
	).Scan(&a.userAlias, &a.password); err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}

		return err
	}

	err := bcrypt.CompareHashAndPassword(a.password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidPassword
	}
	return err
}
