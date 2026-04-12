package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// generateSecureOTP returns a cryptographically secure 6-digit string.
func GenerateSecureOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// DetectType checks if a string is a valid email or mobile number (E.164 format)
func DetectIdentifierType(input string) string {
	input = strings.TrimSpace(input)

	// Check for Email using the standard library
	if _, err := mail.ParseAddress(input); err == nil {
		// Standard ParseAddress can accept "Name <email@domain.com>",
		// so we check if the input contains '@' to be safe for simple strings.
		if strings.Contains(input, "@") && !strings.Contains(input, " ") {
			return "email"
		}
	}

	// Check for Mobile Number (E.164 format: +1234567890)
	// This regex matches an optional '+' followed by 10 to 15 digits.
	mobileRegex := regexp.MustCompile(`^\+?[1-9]\d{9,14}$`)
	if mobileRegex.MatchString(input) {
		return "mobile"
	}

	return "invalid"
}

func HashPassword(password string) ([]byte, error) {

	if password == "" {
		return []byte{}, errors.New("password required")
	}

	bytePassword := []byte(password)
	return bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
}

func CompareHashAndPassword(passwordHash string, password string) error {

	if passwordHash == "" || password == "" {
		return errors.New("password & passwordHash required")
	}

	byteHashedPassword := []byte(passwordHash)
	bytePassword := []byte(password)

	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func RandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes)[:length], nil
}
