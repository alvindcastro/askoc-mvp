package banner

import (
	"net/http"
	"strings"

	"askoc-mvp/internal/fixtures"
	"askoc-mvp/internal/handlers"
)

type Handler struct {
	fixture *fixtures.Fixture
}

type StudentProfile struct {
	StudentID        string   `json:"student_id"`
	PreferredName    string   `json:"preferred_name"`
	Status           string   `json:"status"`
	EnrollmentStatus string   `json:"enrollment_status"`
	Program          string   `json:"program"`
	Holds            []string `json:"holds"`
	Synthetic        bool     `json:"synthetic"`
}

type TranscriptStatus struct {
	StudentID               string   `json:"student_id"`
	TranscriptRequestID     string   `json:"transcript_request_id"`
	TranscriptRequestStatus string   `json:"transcript_request_status"`
	EligibleForProcessing   bool     `json:"eligible_for_processing"`
	Hold                    string   `json:"hold"`
	Holds                   []string `json:"holds"`
	DeliveryMethod          *string  `json:"delivery_method,omitempty"`
	RequestedAt             *string  `json:"requested_at,omitempty"`
	Synthetic               bool     `json:"synthetic"`
}

func NewHandler(fixture *fixtures.Fixture) http.Handler {
	return &Handler{fixture: fixture}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	studentID, suffix, ok := parseStudentPath(r.URL.Path)
	if !ok {
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
		return
	}
	student, found := h.fixture.StudentByID(studentID)
	if !found {
		handlers.WriteError(w, r, http.StatusNotFound, "student_not_found", "No synthetic student record was found for that ID.")
		return
	}

	switch suffix {
	case "":
		handlers.WriteJSON(w, r, http.StatusOK, profileFrom(student))
	case "transcript-status":
		handlers.WriteJSON(w, r, http.StatusOK, transcriptFrom(student))
	default:
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
	}
}

func profileFrom(student fixtures.StudentRecord) StudentProfile {
	return StudentProfile{
		StudentID:        student.StudentID,
		PreferredName:    student.PreferredName,
		Status:           student.Status,
		EnrollmentStatus: student.Status,
		Program:          student.Program,
		Holds:            append([]string(nil), student.Transcript.Holds...),
		Synthetic:        true,
	}
}

func transcriptFrom(student fixtures.StudentRecord) TranscriptStatus {
	hold := "none"
	eligible := true
	for _, value := range student.Transcript.Holds {
		switch value {
		case "mock_financial_hold":
			hold = "financial"
			eligible = false
		case "mock_payment_hold":
			if hold == "none" {
				hold = "payment"
			}
			eligible = false
		}
	}
	if student.Transcript.Status == "not_found" || student.Transcript.Status == "needs_staff_review" {
		eligible = false
	}

	return TranscriptStatus{
		StudentID:               student.StudentID,
		TranscriptRequestID:     student.Transcript.RequestID,
		TranscriptRequestStatus: student.Transcript.Status,
		EligibleForProcessing:   eligible,
		Hold:                    hold,
		Holds:                   append([]string(nil), student.Transcript.Holds...),
		DeliveryMethod:          student.Transcript.DeliveryMethod,
		RequestedAt:             student.Transcript.RequestedAt,
		Synthetic:               true,
	}
}

func parseStudentPath(path string) (studentID string, suffix string, ok bool) {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 4 || len(parts) > 5 {
		return "", "", false
	}
	if parts[0] != "api" || parts[1] != "v1" || parts[2] != "students" || parts[3] == "" {
		return "", "", false
	}
	if len(parts) == 5 {
		suffix = parts[4]
	}
	return parts[3], suffix, true
}
