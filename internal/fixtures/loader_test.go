package fixtures

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFixtureLoadsExpectedSyntheticStudents(t *testing.T) {
	got, err := Load(context.Background(), "../../data/synthetic-students.json")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	for _, id := range []string{"S100001", "S100002", "S100003", "S100004"} {
		if _, ok := got.StudentByID(id); !ok {
			t.Fatalf("StudentByID(%q) missing from fixture", id)
		}
	}

	paid, _ := got.StudentByID("S100001")
	if paid.Payment.Status != "paid" || paid.Transcript.Status != "ready_for_processing" {
		t.Fatalf("S100001 state = payment %q transcript %q", paid.Payment.Status, paid.Transcript.Status)
	}
	if paid.LMS.AccountStatus == "" {
		t.Fatalf("S100001 LMS account status is empty")
	}

	unpaid, _ := got.StudentByID("S100002")
	if unpaid.Payment.Status != "unpaid" || stringValue(unpaid.Payment.TransactionID) != "SYNTH-PAY-100002" {
		t.Fatalf("S100002 payment = %+v", unpaid.Payment)
	}

	hold, _ := got.StudentByID("S100003")
	if !contains(hold.Transcript.Holds, "mock_financial_hold") {
		t.Fatalf("S100003 holds = %v, want mock_financial_hold", hold.Transcript.Holds)
	}

	fallback, _ := got.StudentByID("S100004")
	if fallback.Transcript.Status != "not_found" || fallback.Payment.Status != "not_applicable" {
		t.Fatalf("S100004 state = transcript %q payment %q", fallback.Transcript.Status, fallback.Payment.Status)
	}
}

func TestLoadFixtureRejectsDuplicateStudentIDs(t *testing.T) {
	path := writeFixture(t, `{
		"fixture_name": "duplicate-demo",
		"fixture_version": "test",
		"synthetic_only": true,
		"students": [
			`+studentJSON("S100001", "Demo Learner One")+`,
			`+studentJSON("S100001", "Demo Learner Duplicate")+`
		]
	}`)

	_, err := Load(context.Background(), path)
	if err == nil || !strings.Contains(err.Error(), "duplicate student_id") {
		t.Fatalf("Load() error = %v, want duplicate student_id", err)
	}
}

func TestLoadFixtureRejectsMissingRequiredFields(t *testing.T) {
	path := writeFixture(t, `{
		"fixture_name": "missing-field-demo",
		"fixture_version": "test",
		"synthetic_only": true,
		"students": [
			{
				"student_id": "S100001",
				"synthetic": true,
				"demo_record": true,
				"status": "active",
				"transcript": {"request_id": "SYNTH-TRN-100001", "status": "ready_for_processing", "holds": []},
				"payment": {"status": "paid", "amount_due": 0, "currency": "CAD", "transaction_id": "SYNTH-PAY-100001"},
				"crm": {"case_status": "none"},
				"lms": {"account_status": "active", "courses": []}
			}
		]
	}`)

	_, err := Load(context.Background(), path)
	if err == nil || !strings.Contains(err.Error(), "preferred_name") {
		t.Fatalf("Load() error = %v, want preferred_name validation", err)
	}
}

func TestLoadFixtureRejectsNonSyntheticRecords(t *testing.T) {
	path := writeFixture(t, `{
		"fixture_name": "non-synthetic-demo",
		"fixture_version": "test",
		"synthetic_only": true,
		"students": [
			`+strings.Replace(studentJSON("A100001", "Demo Learner Bad ID"), `"synthetic": true`, `"synthetic": false`, 1)+`
		]
	}`)

	_, err := Load(context.Background(), path)
	if err == nil || !strings.Contains(err.Error(), "synthetic") {
		t.Fatalf("Load() error = %v, want synthetic validation", err)
	}
}

func writeFixture(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "students.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func studentJSON(id, name string) string {
	return `{
		"student_id": "` + id + `",
		"synthetic": true,
		"demo_record": true,
		"preferred_name": "` + name + `",
		"program": "Demo Program",
		"status": "active",
		"transcript": {"request_id": "SYNTH-TRN-100001", "status": "ready_for_processing", "holds": []},
		"payment": {"status": "paid", "amount_due": 0, "currency": "CAD", "transaction_id": "SYNTH-PAY-100001"},
		"crm": {"case_status": "none"},
		"lms": {"account_status": "active", "courses": [{"course_id": "DEMO-LMS-101", "course_name": "Demo Course", "access_status": "available"}]}
	}`
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
