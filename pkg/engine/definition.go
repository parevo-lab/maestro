package engine

import (
	"time"
)

// WorkflowDefinition bir iş akışının tanımını temsil eder
type WorkflowDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     int                    `json:"version"`
	Steps       []StepDefinition       `json:"steps"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// StepDefinition bir iş akışı adımının tanımını temsil eder
type StepDefinition struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        StepType               `json:"type"`
	Config      map[string]interface{} `json:"config,omitempty"`
	NextSteps   []string               `json:"next_steps,omitempty"`
	RetryPolicy *RetryPolicy           `json:"retry_policy,omitempty"`
	Timeout     time.Duration          `json:"timeout,omitempty"`
}

// StepType adım tiplerini temsil eder
type StepType string

const (
	StepTypeTask     StepType = "task"
	StepTypeApproval StepType = "approval"
	StepTypeDecision StepType = "decision"
	StepTypeProcess  StepType = "process"
)

// RetryPolicy yeniden deneme politikasını temsil eder
type RetryPolicy struct {
	MaxAttempts     int           `json:"max_attempts"`
	InitialInterval time.Duration `json:"initial_interval"`
	MaxInterval     time.Duration `json:"max_interval"`
	Multiplier      float64       `json:"multiplier"`
}

// NewWorkflowDefinition yeni bir iş akışı tanımı oluşturur
func NewWorkflowDefinition(id, name, description string) *WorkflowDefinition {
	now := time.Now()
	return &WorkflowDefinition{
		ID:          id,
		Name:        name,
		Description: description,
		Version:     1,
		Steps:       make([]StepDefinition, 0),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddStep iş akışına yeni bir adım ekler
func (w *WorkflowDefinition) AddStep(step StepDefinition) {
	w.Steps = append(w.Steps, step)
	w.UpdatedAt = time.Now()
}

// NewStepDefinition yeni bir adım tanımı oluşturur
func NewStepDefinition(id, name string, stepType StepType) StepDefinition {
	return StepDefinition{
		ID:        id,
		Name:      name,
		Type:      stepType,
		Config:    make(map[string]interface{}),
		NextSteps: make([]string, 0),
	}
}

// WithConfig adıma yapılandırma ekler
func (s StepDefinition) WithConfig(config map[string]interface{}) StepDefinition {
	s.Config = config
	return s
}

// WithNextSteps adıma sonraki adımları ekler
func (s StepDefinition) WithNextSteps(nextSteps ...string) StepDefinition {
	s.NextSteps = nextSteps
	return s
}

// WithRetryPolicy adıma yeniden deneme politikası ekler
func (s StepDefinition) WithRetryPolicy(maxAttempts int, initialInterval, maxInterval time.Duration, multiplier float64) StepDefinition {
	s.RetryPolicy = &RetryPolicy{
		MaxAttempts:     maxAttempts,
		InitialInterval: initialInterval,
		MaxInterval:     maxInterval,
		Multiplier:      multiplier,
	}
	return s
}

// WithTimeout adıma zaman aşımı süresi ekler
func (s StepDefinition) WithTimeout(timeout time.Duration) StepDefinition {
	s.Timeout = timeout
	return s
}
