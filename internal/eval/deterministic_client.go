package eval

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/fixtures"
	"askoc-mvp/internal/orchestrator"
	"askoc-mvp/internal/rag"
	"askoc-mvp/internal/tools"
	"askoc-mvp/internal/workflow"
)

type DeterministicClientConfig struct {
	ChunksPath   string
	StudentsPath string
}

type serviceChatClient struct {
	service interface {
		HandleChat(context.Context, domain.ChatRequest) (domain.ChatResponse, error)
	}
}

func (c serviceChatClient) Chat(ctx context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	return c.service.HandleChat(ctx, req)
}

func NewDeterministicChatClient(ctx context.Context, cfg DeterministicClientConfig) (ChatClient, error) {
	if strings.TrimSpace(cfg.ChunksPath) == "" {
		cfg.ChunksPath = "data/rag-chunks.json"
	}
	if strings.TrimSpace(cfg.StudentsPath) == "" {
		cfg.StudentsPath = "data/synthetic-students.json"
	}

	chunks, err := rag.LoadChunks(ctx, cfg.ChunksPath)
	if err != nil {
		return nil, err
	}
	fixture, err := fixtures.Load(ctx, cfg.StudentsPath)
	if err != nil {
		return nil, err
	}
	crm := &fixtureCRM{}
	service, err := orchestrator.New(orchestrator.Dependencies{
		Classifier: classifier.Fallback{},
		Retriever:  rag.NewLocalRetriever(chunks),
		LLM:        orchestrator.DisabledLLM{},
		Banner:     fixtureBanner{fixture: fixture},
		Payment:    fixturePayment{fixture: fixture},
		Workflow:   workflow.NewInMemoryClient(),
		CRM:        crm,
		Audit:      audit.NewMemoryStore(),
	})
	if err != nil {
		return nil, err
	}
	return serviceChatClient{service: service}, nil
}

type fixtureBanner struct {
	fixture *fixtures.Fixture
}

func (b fixtureBanner) GetTranscriptStatus(ctx context.Context, studentID string) (tools.BannerTranscriptStatus, error) {
	if err := ctx.Err(); err != nil {
		return tools.BannerTranscriptStatus{}, err
	}
	student, ok := b.fixture.StudentByID(strings.TrimSpace(strings.ToUpper(studentID)))
	if !ok {
		return tools.BannerTranscriptStatus{}, errors.New("synthetic transcript record not found")
	}
	hold := ""
	if len(student.Transcript.Holds) > 0 {
		hold = student.Transcript.Holds[0]
	}
	return tools.BannerTranscriptStatus{
		StudentID:               student.StudentID,
		TranscriptRequestID:     student.Transcript.RequestID,
		TranscriptRequestStatus: student.Transcript.Status,
		EligibleForProcessing:   len(student.Transcript.Holds) == 0 && student.Transcript.Status != "not_found",
		Hold:                    hold,
		Holds:                   append([]string(nil), student.Transcript.Holds...),
		Synthetic:               true,
	}, nil
}

type fixturePayment struct {
	fixture *fixtures.Fixture
}

func (p fixturePayment) GetPaymentStatus(ctx context.Context, studentID string) (tools.PaymentStatus, error) {
	if err := ctx.Err(); err != nil {
		return tools.PaymentStatus{}, err
	}
	student, ok := p.fixture.StudentByID(strings.TrimSpace(strings.ToUpper(studentID)))
	if !ok {
		return tools.PaymentStatus{}, errors.New("synthetic payment record not found")
	}
	transactionID := ""
	if student.Payment.TransactionID != nil {
		transactionID = *student.Payment.TransactionID
	}
	return tools.PaymentStatus{
		StudentID:     student.StudentID,
		Item:          "official_transcript",
		AmountDue:     student.Payment.AmountDue,
		Currency:      student.Payment.Currency,
		Status:        student.Payment.Status,
		TransactionID: transactionID,
		Synthetic:     true,
	}, nil
}

type fixtureCRM struct{}

func (c *fixtureCRM) CreateCase(ctx context.Context, req tools.CRMCaseRequest) (tools.CRMCaseResponse, error) {
	if err := ctx.Err(); err != nil {
		return tools.CRMCaseResponse{}, err
	}
	if strings.TrimSpace(req.Queue) == "" {
		return tools.CRMCaseResponse{}, errors.New("crm queue is required")
	}
	return tools.CRMCaseResponse{
		CaseID:         deterministicCaseID(req),
		Status:         "queued",
		Queue:          req.Queue,
		Priority:       req.Priority,
		Summary:        req.Summary,
		ConversationID: req.ConversationID,
		SourceTraceID:  req.SourceTraceID,
		Synthetic:      true,
	}, nil
}

func deterministicCaseID(req tools.CRMCaseRequest) string {
	key := fmt.Sprintf("%s:%s:%s:%s", req.StudentID, req.ConversationID, req.Queue, req.SourceTraceID)
	hash := sha1.Sum([]byte(key))
	return "EVAL-CRM-" + strings.ToUpper(hex.EncodeToString(hash[:])[:10])
}
