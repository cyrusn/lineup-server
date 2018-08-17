package auth

import (
	"database/sql"
	"errors"
	"time"

	auth "github.com/cyrusn/goJWTAuthHelper"
	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")

	// lifeTime is jwt life Time in minutes
	lifeTime int64
)

// UpdateLifeTime is used for update the lifeTime of jwt token
func UpdateLifeTime(time int64) {
	lifeTime = time
}

// Credential store user credential
type Credential struct {
	UserAlias string `json:"userAlias"`
	Password  string `json:"password"`
	Role      string `json:role`
}

// DB is *sql.DB store auth information
type DB struct {
	*sql.DB
	*auth.Secret
}

// Claims is jwt.Claims for authentication credential
type Claims struct {
	UserAlias string
	Role      string
	jwt.StandardClaims
}

// Update use to refresh token,
func (claims *Claims) Update(token *jwt.Token) {
	claims.ExpiresAt = expiresAfter(lifeTime)
	token.Claims = claims
}

// Insert insert credential for authentication
func (db *DB) Insert(userAlias, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO Credential (
      useralias,
      password,
      role
    ) values (?, ?, ?)`,
		userAlias, hashedPassword, role,
	)
	return err
}

// Authenticate will lookup userAlias in sqlite3 database and
// compare the hashed password store in database
// return jwtToken and error
func (db *DB) Authenticate(userAlias, password string) (string, error) {
	type credential struct {
		userAlias string
		password  []byte
		role      string
	}

	c := new(credential)

	err := db.QueryRow(
		`SELECT * FROM Credential WHERE useralias=?`,
		userAlias,
	).Scan(&c.userAlias, &c.password, &c.role)

	switch {
	case err == sql.ErrNoRows:
		return "", ErrUserNotFound
	case err != nil:
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(c.password, []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", ErrInvalidPassword
		}
		return "", err
	}

	claims := Claims{
		UserAlias: userAlias,
		Role:      c.role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAfter(lifeTime),
		},
	}
	return db.Secret.GenerateToken(claims)
}

// Refresh refresh the jwt token
func (db *DB) Refresh(jwt string) (string, error) {
	claims := new(Claims)
	return db.Secret.UpdateToken(jwt, claims)
}

func expiresAfter(min int64) int64 {
	return time.Now().Add(time.Minute * time.Duration(min)).Unix()
}
