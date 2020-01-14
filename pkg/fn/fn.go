package fn

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ErrInvalidEmail = errors.New("invalid email address")
var userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
var hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")

func init() { rand.Seed(time.Now().UnixNano()) }

// GenerateRandomString returns a randomly generated string
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidEmail
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return ErrInvalidEmail
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return ErrInvalidEmail
	}

	if !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return ErrInvalidEmail
	}

	return nil
}

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func VerifyHashPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	}
	return true
}
