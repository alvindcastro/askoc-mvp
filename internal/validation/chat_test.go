package validation

import (
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
)

func TestValidateChatRequestMessageRules(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		code string
	}{
		{name: "empty", msg: "", code: CodeMissingMessage},
		{name: "whitespace", msg: " \n\t ", code: CodeMissingMessage},
		{name: "oversized", msg: strings.Repeat("a", MaxChatMessageLength+1), code: CodeMessageTooLarge},
		{name: "valid at limit", msg: strings.Repeat("a", MaxChatMessageLength), code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChatRequest(domain.ChatRequest{Channel: "web", Message: tt.msg})
			if tt.code == "" {
				if err != nil {
					t.Fatalf("ValidateChatRequest error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatal("ValidateChatRequest error = nil, want validation error")
			}
			if code := Code(err); code != tt.code {
				t.Fatalf("validation code = %q, want %q", code, tt.code)
			}
		})
	}
}

func TestValidateChatRequestStudentIDRules(t *testing.T) {
	tests := []struct {
		name      string
		studentID string
		wantErr   bool
	}{
		{name: "empty optional", studentID: "", wantErr: false},
		{name: "synthetic ID", studentID: "S100002", wantErr: false},
		{name: "synthetic not found shape still valid", studentID: "S999999", wantErr: false},
		{name: "missing prefix", studentID: "100002", wantErr: true},
		{name: "lowercase prefix", studentID: "s100002", wantErr: true},
		{name: "too short", studentID: "S10002", wantErr: true},
		{name: "non-digit suffix", studentID: "S10000X", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChatRequest(domain.ChatRequest{
				Channel:   "web",
				Message:   "Please check my transcript.",
				StudentID: tt.studentID,
			})
			if tt.wantErr && err == nil {
				t.Fatal("ValidateChatRequest error = nil, want validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateChatRequest error = %v, want nil", err)
			}
			if tt.wantErr && Code(err) != CodeInvalidStudentID {
				t.Fatalf("validation code = %q, want %q", Code(err), CodeInvalidStudentID)
			}
		})
	}
}
