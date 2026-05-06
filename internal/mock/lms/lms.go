package lms

import (
	"net/http"
	"strings"

	"askoc-mvp/internal/fixtures"
	"askoc-mvp/internal/handlers"
)

type Handler struct {
	fixture *fixtures.Fixture
}

type AccessStatus struct {
	StudentID       string `json:"student_id"`
	AccountStatus   string `json:"account_status"`
	CourseID        string `json:"course_id,omitempty"`
	CourseName      string `json:"course_name,omitempty"`
	AccessStatus    string `json:"access_status"`
	Message         string `json:"message,omitempty"`
	Synthetic       bool   `json:"synthetic"`
	ContentIncluded bool   `json:"content_included"`
}

func NewHandler(fixture *fixtures.Fixture) http.Handler {
	return &Handler{fixture: fixture}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	studentID, ok := parseLMSPath(r.URL.Path)
	if !ok {
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
		return
	}
	student, found := h.fixture.StudentByID(studentID)
	if !found {
		handlers.WriteError(w, r, http.StatusNotFound, "lms_student_not_found", "No synthetic LMS record was found for that ID.")
		return
	}

	courseID := strings.TrimSpace(r.URL.Query().Get("course_id"))
	for _, course := range student.LMS.Courses {
		if course.CourseID == courseID {
			handlers.WriteJSON(w, r, http.StatusOK, AccessStatus{
				StudentID:       student.StudentID,
				AccountStatus:   student.LMS.AccountStatus,
				CourseID:        course.CourseID,
				CourseName:      course.CourseName,
				AccessStatus:    course.AccessStatus,
				Synthetic:       true,
				ContentIncluded: false,
			})
			return
		}
	}

	handlers.WriteJSON(w, r, http.StatusOK, AccessStatus{
		StudentID:       student.StudentID,
		AccountStatus:   student.LMS.AccountStatus,
		CourseID:        courseID,
		AccessStatus:    "unknown_demo_course",
		Message:         "No synthetic LMS access record exists for that demo course.",
		Synthetic:       true,
		ContentIncluded: false,
	})
}

func parseLMSPath(path string) (string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 5 {
		return "", false
	}
	if parts[0] != "api" || parts[1] != "v1" || parts[2] != "students" || parts[3] == "" || parts[4] != "lms-access" {
		return "", false
	}
	return parts[3], true
}
