package tools

import (
	"context"
	"net/http"
	"strings"
)

type CRMClient struct {
	baseURL    string
	httpClient *http.Client
}

type CRMCaseRequest struct {
	StudentID      string `json:"student_id"`
	ConversationID string `json:"conversation_id"`
	Intent         string `json:"intent"`
	Priority       string `json:"priority"`
	Queue          string `json:"queue"`
	Summary        string `json:"summary"`
	SourceTraceID  string `json:"source_trace_id"`
}

type CRMCaseResponse struct {
	CaseID         string `json:"case_id"`
	Status         string `json:"status"`
	Queue          string `json:"queue"`
	Priority       string `json:"priority"`
	Summary        string `json:"summary"`
	ConversationID string `json:"conversation_id,omitempty"`
	SourceTraceID  string `json:"source_trace_id,omitempty"`
	Synthetic      bool   `json:"synthetic"`
}

func NewCRMClient(baseURL string, httpClient *http.Client) *CRMClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &CRMClient{baseURL: strings.TrimRight(baseURL, "/"), httpClient: httpClient}
}

func (c *CRMClient) CreateCase(ctx context.Context, req CRMCaseRequest) (CRMCaseResponse, error) {
	var out CRMCaseResponse
	err := doJSON(ctx, c.httpClient, "crm", http.MethodPost, c.baseURL+"/api/v1/crm/cases", req, &out)
	return out, err
}
