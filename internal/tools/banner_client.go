package tools

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type BannerClient struct {
	baseURL    string
	httpClient *http.Client
}

type BannerStudentProfile struct {
	StudentID        string   `json:"student_id"`
	PreferredName    string   `json:"preferred_name"`
	Status           string   `json:"status"`
	EnrollmentStatus string   `json:"enrollment_status"`
	Program          string   `json:"program"`
	Holds            []string `json:"holds"`
	Synthetic        bool     `json:"synthetic"`
}

type BannerTranscriptStatus struct {
	StudentID               string   `json:"student_id"`
	TranscriptRequestID     string   `json:"transcript_request_id"`
	TranscriptRequestStatus string   `json:"transcript_request_status"`
	EligibleForProcessing   bool     `json:"eligible_for_processing"`
	Hold                    string   `json:"hold"`
	Holds                   []string `json:"holds"`
	Synthetic               bool     `json:"synthetic"`
}

func NewBannerClient(baseURL string, httpClient *http.Client) *BannerClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &BannerClient{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

func (c *BannerClient) GetStudent(ctx context.Context, studentID string) (BannerStudentProfile, error) {
	var out BannerStudentProfile
	err := doJSON(ctx, c.httpClient, "banner", http.MethodGet, c.baseURL+"/api/v1/students/"+url.PathEscape(studentID), nil, &out)
	return out, err
}

func (c *BannerClient) GetTranscriptStatus(ctx context.Context, studentID string) (BannerTranscriptStatus, error) {
	var out BannerTranscriptStatus
	err := doJSON(ctx, c.httpClient, "banner", http.MethodGet, c.baseURL+"/api/v1/students/"+url.PathEscape(studentID)+"/transcript-status", nil, &out)
	return out, err
}
