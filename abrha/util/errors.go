package util

import (
	"strings"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
)

// IsAbrhaError detects if a given error is a *goApiAbrha.ErrorResponse for
// the specified code and message.
func IsAbrhaError(err error, code int, message string) bool {
	if err, ok := err.(*goApiAbrha.ErrorResponse); ok {
		return err.Response.StatusCode == code &&
			strings.Contains(strings.ToLower(err.Message), strings.ToLower(message))
	}
	return false
}
