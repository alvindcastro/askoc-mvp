package tools

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	KindNotFound  ErrorKind = "not_found"
	KindRetryable ErrorKind = "retryable"
	KindExternal  ErrorKind = "external_service"
	KindParse     ErrorKind = "parse"
	KindTimeout   ErrorKind = "timeout"
)

type ToolError struct {
	Kind       ErrorKind
	Service    string
	StatusCode int
	Message    string
	Err        error
}

func (e *ToolError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return fmt.Sprintf("%s %s: %s", e.Service, e.Kind, e.Message)
	}
	return fmt.Sprintf("%s %s", e.Service, e.Kind)
}

func (e *ToolError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func IsKind(err error, kind ErrorKind) bool {
	var toolErr *ToolError
	return errors.As(err, &toolErr) && toolErr.Kind == kind
}
