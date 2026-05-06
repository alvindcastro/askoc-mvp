package tools

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type LMSClient struct {
	baseURL    string
	httpClient *http.Client
}

type LMSAccessStatus struct {
	StudentID       string `json:"student_id"`
	AccountStatus   string `json:"account_status"`
	CourseID        string `json:"course_id,omitempty"`
	CourseName      string `json:"course_name,omitempty"`
	AccessStatus    string `json:"access_status"`
	Message         string `json:"message,omitempty"`
	Synthetic       bool   `json:"synthetic"`
	ContentIncluded bool   `json:"content_included"`
}

func NewLMSClient(baseURL string, httpClient *http.Client) *LMSClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &LMSClient{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

func (c *LMSClient) GetAccessStatus(ctx context.Context, studentID, courseID string) (LMSAccessStatus, error) {
	endpoint := c.baseURL + "/api/v1/students/" + url.PathEscape(studentID) + "/lms-access"
	if strings.TrimSpace(courseID) != "" {
		values := url.Values{}
		values.Set("course_id", courseID)
		endpoint += "?" + values.Encode()
	}
	var out LMSAccessStatus
	err := doJSON(ctx, c.httpClient, "lms", http.MethodGet, endpoint, nil, &out)
	return out, err
}
