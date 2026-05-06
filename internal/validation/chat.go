package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"askoc-mvp/internal/domain"
)

const (
	MaxChatMessageLength = 2000

	CodeInvalidRequest   = "invalid_request"
	CodeMissingMessage   = "missing_message"
	CodeMessageTooLarge  = "message_too_large"
	CodeInvalidStudentID = "invalid_student_id"
)

var syntheticStudentID = regexp.MustCompile(`^S[0-9]{6}$`)

type ValidationError struct {
	code    string
	message string
}

func (e ValidationError) Error() string {
	return e.message
}

func Code(err error) string {
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.code
	}
	return CodeInvalidRequest
}

func SafeMessage(err error) string {
	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.message
	}
	return "request is invalid"
}

func ValidateChatRequest(req domain.ChatRequest) error {
	message := strings.TrimSpace(req.Message)
	if message == "" {
		return ValidationError{
			code:    CodeMissingMessage,
			message: "message is required",
		}
	}
	if utf8.RuneCountInString(req.Message) > MaxChatMessageLength {
		return ValidationError{
			code:    CodeMessageTooLarge,
			message: "message must be 2000 characters or fewer",
		}
	}
	if strings.TrimSpace(req.StudentID) != "" && !syntheticStudentID.MatchString(req.StudentID) {
		return ValidationError{
			code:    CodeInvalidStudentID,
			message: "student_id must use the synthetic S plus six digits demo format",
		}
	}
	return nil
}
