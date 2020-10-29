package util

import (
	"bytes"
	"fmt"
	"github.com/go-playground/validator"
	"neoway-case/errors"
	"neoway-case/schema"
	"net/http"
	"regexp"
	"strconv"
	"unicode"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation(`cpf`, isCPF, true)
	validate.RegisterValidation(`cnpj`, isCNPJ, true)
}

//Validate every field from Consumption struct type. First error encountered is sent back to caller.
func Validate(consumptions []schema.Consumption) error {
	for i, c := range consumptions {
		if err := validate.Struct(c); err != nil {
			message := "Validation error at line " + strconv.Itoa(i+2) + ". "
			message = getErrorText(err, message)
			return errors.E(errors.Op("util.Validate"), errors.Message(message), http.StatusBadRequest, errors.Kind("Validation error"))
		}
	}
	return nil
}

//Get a string by concatening field value and which validation failed.
func getErrorText(err error, message string) string {
	for _, err := range err.(validator.ValidationErrors) {
		message += "Field " + err.Field() + " with value " + fmt.Sprintf("%v", err.Value()) + " failed " + err.Tag() + " tag validation. "
	}
	return message
}

// Regexp pattern for CPF and CNPJ.
var (
	CPFRegexp  = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)
	CNPJRegexp = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)
)

// isCPF verifies if the given string is a valid CPF document.
func isCPF(f1 validator.FieldLevel) bool {

	const (
		size = 9
		pos  = 10
	)

	return isCPFOrCNPJ(f1.Field().String(), CPFRegexp, size, pos)
}

// isCNPJ verifies if the given string is a valid CNPJ document.
func isCNPJ(f1 validator.FieldLevel) bool {

	const (
		size = 12
		pos  = 5
	)
	cnpj := f1.Field().String()
	if cnpj == "<*string Value>" {
		return true
	}
	return isCPFOrCNPJ(cnpj, CNPJRegexp, size, pos)
}

// isCPFOrCNPJ generates the digits for a given CPF or CNPJ and compares it with the original digits.
func isCPFOrCNPJ(doc string, pattern *regexp.Regexp, size int, position int) bool {

	if !pattern.MatchString(doc) {
		return false
	}

	cleanNonDigits(&doc)

	// Invalidates documents with all digits equal.
	if allEq(doc) {
		return false
	}

	d := doc[:size]
	digit := calculateDigit(d, position)

	d = d + digit
	digit = calculateDigit(d, position+1)

	return doc == d+digit
}

// cleanNonDigits removes every rune that is not a digit.
func cleanNonDigits(doc *string) {

	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}

// allEq checks if every rune in a given string is equal.
func allEq(doc string) bool {

	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

// calculateDigit calculates the next digit for the given document.
func calculateDigit(doc string, position int) string {

	var sum int
	for _, r := range doc {

		sum += toInt(r) * position
		position--

		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}

	return strconv.Itoa(11 - sum)
}

func toInt(r rune) int {
	return int(r - '0')
}
