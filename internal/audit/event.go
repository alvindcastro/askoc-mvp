package audit

import "context"

type Event struct {
	TraceID        string
	ConversationID string
	StudentID      string
	Type           string
	Action         string
	Status         string
	ReferenceID    string
	Message        string
	Metadata       map[string]string
}

type Recorder interface {
	Record(context.Context, Event) error
}

type NopRecorder struct{}

func (NopRecorder) Record(context.Context, Event) error {
	return nil
}
