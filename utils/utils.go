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

// create a cryptographically secure 6-digit string.
func GenerateSecureOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// checks if a string is a valid email or mobile number or otp-resend-key
func DetectIdentifierType(input string) string {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "otp-resend-") {
		return "otp-resend"
	}

	// Email Check using the standard library
	if _, err := mail.ParseAddress(input); err == nil {
		// does input contains '@' 
		if strings.Contains(input, "@") && !strings.Contains(input, " ") {
			return "email"
		}
	}

	// Check for Mobile Number (matches an optional '+' followed by 10 to 15 digits.)
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
