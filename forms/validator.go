package forms

import (
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/KibuuleNoah/QuickGin/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// DefaultValidator ...
type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

// ValidateStruct ...
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

// Engine ...
func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {

		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here

		//Custom rule for user full name
		v.validate.RegisterValidation("username", ValidateFullName)
		v.validate.RegisterValidation("identifier", ValidateIdentifier)
		v.validate.RegisterValidation("strong_password", ValidateStrongPassword)
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValidateFullName implements validator.Func
func ValidateFullName(fl validator.FieldLevel) bool {
	//Remove the extra space
	space := regexp.MustCompile(`\s+`)
	name := space.ReplaceAllString(fl.Field().String(), " ")

	//Remove trailing spaces
	name = strings.TrimSpace(name)

	//To support all possible languages
	matched, _ := regexp.Match(`^[^±!@£$%^&*_+§¡€#¢§¶•ªº«\\/<>?:;'"|=.,0123456789]{3,20}$`, []byte(name))
	return matched
}

func ValidateIdentifier(fl validator.FieldLevel) bool {
	return utils.DetectIdentifierType(fl.Field().String()) != "invalid"
}

const (
	minPasswordLength = 8
	maxPasswordLength = 32

	minUppercase = 1
	minLowercase = 1
	minDigits    = 1
	minSpecial   = 1
)

// passwordStrength holds per-character-class counts after a single scan.
type passwordStrength struct {
	length       int
	uppercase    int
	lowercase    int
	digits       int
	special      int
	hasAmbiguous bool // optional: flag visually-ambiguous chars (0/O, 1/l/I)
}

// analyzePassword performs a single O(n) scan of the password rune slice.
func analyzePassword(password string) passwordStrength {
	var s passwordStrength
	for _, r := range password {
		s.length++
		switch {
		case unicode.IsUpper(r):
			s.uppercase++
		case unicode.IsLower(r):
			s.lowercase++
		case unicode.IsDigit(r):
			s.digits++
		case isSpecialChar(r):
			s.special++
		}
		if r == '0' || r == 'O' || r == '1' || r == 'l' || r == 'I' {
			s.hasAmbiguous = true
		}
	}
	return s
}

// isSpecialChar returns true for printable, non-alphanumeric ASCII characters
// plus common Unicode punctuation/symbols.
func isSpecialChar(r rune) bool {
	if r > unicode.MaxASCII {
		return unicode.IsPunct(r) || unicode.IsSymbol(r)
	}
	// Printable ASCII special characters (33–126, excluding letters and digits)
	return r >= '!' && r <= '~' && !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

// hasSequentialChars detects runs of 3+ sequential characters (e.g. "abc", "123").
func hasSequentialChars(password string) bool {
	runes := []rune(password)
	if len(runes) < 3 {
		return false
	}
	for i := 0; i < len(runes)-2; i++ {
		diff1 := runes[i+1] - runes[i]
		diff2 := runes[i+2] - runes[i+1]
		if diff1 == 1 && diff2 == 1 {
			return true
		}
		if diff1 == -1 && diff2 == -1 {
			return true
		}
	}
	return false
}

// hasRepeatingChars detects 3+ consecutive identical characters (e.g. "aaa").
func hasRepeatingChars(password string) bool {
	runes := []rune(password)
	count := 1
	for i := 1; i < len(runes); i++ {
		if runes[i] == runes[i-1] {
			count++
			if count >= 3 {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

// Password string `binding:"required,strong_password"`
func ValidateStrongPassword(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// password was not supplied
	if password == "" {
		return true
	}

	// Reject invalid UTF-8 early.
	if !utf8.ValidString(password) {
		return false
	}

	// Length bounds (rune-aware, not byte-count).
	runeLen := utf8.RuneCountInString(password)
	if runeLen < minPasswordLength || runeLen > maxPasswordLength {
		return false
	}

	// Sequential / repeating pattern checks.
	if hasSequentialChars(password) {
		return false
	}
	if hasRepeatingChars(password) {
		return false
	}

	// Character-class requirements (single pass).
	s := analyzePassword(password)
	if s.uppercase < minUppercase {
		return false
	}
	if s.lowercase < minLowercase {
		return false
	}
	if s.digits < minDigits {
		return false
	}
	if s.special < minSpecial {
		return false
	}

	return true
}
