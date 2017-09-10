package util

import(
  "fmt"
  "net/http"
)

type StatusError struct {
	Status int
	Message string
}

func (s StatusError) Error() string {
	return s.Message
}

// Get an error with an associated HTTP status code.
func CommandError(status int, message string, a ...interface{}) *StatusError {
  if message == "" {
    message = http.StatusText(status)
  }
	return &StatusError{Status: status, Message: fmt.Sprintf(message, a...)}
}

