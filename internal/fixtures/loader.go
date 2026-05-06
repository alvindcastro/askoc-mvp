package fixtures

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var syntheticStudentIDPattern = regexp.MustCompile(`^S[0-9]{6}$`)

type Fixture struct {
	FixtureName    string          `json:"fixture_name"`
	FixtureVersion string          `json:"fixture_version"`
	SyntheticOnly  bool            `json:"synthetic_only"`
	DataNotice     string          `json:"data_notice"`
	IDPolicy       IDPolicy        `json:"id_policy"`
	Students       []StudentRecord `json:"students"`

	studentsByID map[string]StudentRecord
}

type IDPolicy struct {
	StudentIDPattern    string `json:"student_id_pattern"`
	TransactionIDPrefix string `json:"transaction_id_prefix"`
	CRMCaseIDPrefix     string `json:"crm_case_id_prefix"`
	Notes               string `json:"notes"`
}

type StudentRecord struct {
	StudentID     string           `json:"student_id"`
	Synthetic     bool             `json:"synthetic"`
	DemoRecord    bool             `json:"demo_record"`
	PreferredName string           `json:"preferred_name"`
	Program       string           `json:"program"`
	Status        string           `json:"status"`
	Transcript    TranscriptRecord `json:"transcript"`
	Payment       PaymentRecord    `json:"payment"`
	CRM           CRMRecord        `json:"crm"`
	LMS           LMSRecord        `json:"lms"`
	DemoNotes     string           `json:"demo_notes"`
}

type TranscriptRecord struct {
	RequestID      string   `json:"request_id"`
	Status         string   `json:"status"`
	RequestedAt    *string  `json:"requested_at"`
	DeliveryMethod *string  `json:"delivery_method"`
	Holds          []string `json:"holds"`
}

type PaymentRecord struct {
	Status        string  `json:"status"`
	AmountDue     float64 `json:"amount_due"`
	Currency      string  `json:"currency"`
	TransactionID *string `json:"transaction_id"`
}

type CRMRecord struct {
	CaseID     *string `json:"case_id"`
	CaseStatus string  `json:"case_status"`
}

type LMSRecord struct {
	AccountStatus string      `json:"account_status"`
	Courses       []LMSCourse `json:"courses"`
}

type LMSCourse struct {
	CourseID     string `json:"course_id"`
	CourseName   string `json:"course_name"`
	AccessStatus string `json:"access_status"`
}

func Load(ctx context.Context, path string) (*Fixture, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load synthetic fixture: %w", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load synthetic fixture: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load synthetic fixture: %w", err)
	}

	var fixture Fixture
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		return nil, fmt.Errorf("decode synthetic fixture: %w", err)
	}
	if err := fixture.validate(); err != nil {
		return nil, err
	}
	return &fixture, nil
}

func (f *Fixture) StudentByID(studentID string) (StudentRecord, bool) {
	if f == nil || f.studentsByID == nil {
		return StudentRecord{}, false
	}
	student, ok := f.studentsByID[studentID]
	return student, ok
}

func (f *Fixture) validate() error {
	if !f.SyntheticOnly {
		return fmt.Errorf("synthetic fixture must set synthetic_only=true")
	}
	if len(f.Students) == 0 {
		return fmt.Errorf("synthetic fixture must include students")
	}

	f.studentsByID = make(map[string]StudentRecord, len(f.Students))
	for i, student := range f.Students {
		if err := validateStudent(student); err != nil {
			return fmt.Errorf("student[%d] %s: %w", i, student.StudentID, err)
		}
		if _, exists := f.studentsByID[student.StudentID]; exists {
			return fmt.Errorf("duplicate student_id %q in synthetic fixture", student.StudentID)
		}
		f.studentsByID[student.StudentID] = student
	}
	return nil
}

func validateStudent(student StudentRecord) error {
	if !syntheticStudentIDPattern.MatchString(student.StudentID) {
		return fmt.Errorf("student_id must be synthetic S plus six digits")
	}
	if !student.Synthetic || !student.DemoRecord {
		return fmt.Errorf("record must be synthetic demo data")
	}
	if strings.TrimSpace(student.PreferredName) == "" {
		return fmt.Errorf("preferred_name is required")
	}
	if strings.TrimSpace(student.Program) == "" {
		return fmt.Errorf("program is required")
	}
	if strings.TrimSpace(student.Status) == "" {
		return fmt.Errorf("status is required")
	}
	if strings.TrimSpace(student.Transcript.Status) == "" {
		return fmt.Errorf("transcript.status is required")
	}
	if strings.TrimSpace(student.Payment.Status) == "" {
		return fmt.Errorf("payment.status is required")
	}
	if strings.TrimSpace(student.Payment.Currency) == "" {
		return fmt.Errorf("payment.currency is required")
	}
	if strings.TrimSpace(student.CRM.CaseStatus) == "" {
		return fmt.Errorf("crm.case_status is required")
	}
	if strings.TrimSpace(student.LMS.AccountStatus) == "" {
		return fmt.Errorf("lms.account_status is required")
	}
	for _, course := range student.LMS.Courses {
		if strings.TrimSpace(course.CourseID) == "" || strings.TrimSpace(course.AccessStatus) == "" {
			return fmt.Errorf("lms.courses require course_id and access_status")
		}
	}
	return nil
}
