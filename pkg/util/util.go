package util

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) (string, bool) {
	var (
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
		specialRunes = "!@#$%^&*"
	)

	if len(password) < 8 {
		return "password must be at least 8 characters long", false
	}

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case containsRune(specialRunes, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "password must contain at least one uppercase letter", false
	}
	if !hasLower {
		return "password must contain at least one lowercase letter", false
	}
	if !hasNumber {
		return "password must contain at least one number", false
	}
	if !hasSpecial {
		return "password must contain at least one special character", false
	}

	return "Password is valid", true
}

func containsRune(str string, r rune) bool {
	for _, sr := range str {
		if sr == r {
			return true
		}
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
