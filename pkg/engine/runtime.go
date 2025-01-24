package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkflowRuntime iş akışı çalışma zamanını temsil eder
type WorkflowRuntime struct {
	engine     *WorkflowEngine
	definition *WorkflowDefinition
	state      *WorkflowState
	mutex      sync.RWMutex
}

// WorkflowState iş akışının durumunu temsil eder
type WorkflowState struct {
	CurrentStepID string
	Status        WorkflowStatus
	Context       map[string]interface{}
	StepResults   map[string]interface{}
	StartedAt     time.Time
	CompletedAt   *time.Time
	Error         error
}

// WorkflowStatus iş akışı durumunu temsil eder
type WorkflowStatus string

const (
	StatusPending   WorkflowStatus = "pending"
	StatusRunning   WorkflowStatus = "running"
	StatusCompleted WorkflowStatus = "completed"
	StatusFailed    WorkflowStatus = "failed"
	StatusCanceled  WorkflowStatus = "canceled"
)

// NewWorkflowRuntime yeni bir iş akışı çalışma zamanı oluşturur
func NewWorkflowRuntime(engine *WorkflowEngine, definition *WorkflowDefinition) *WorkflowRuntime {
	return &WorkflowRuntime{
		engine:     engine,
		definition: definition,
		state: &WorkflowState{
			Status:      StatusPending,
			Context:     make(map[string]interface{}),
			StepResults: make(map[string]interface{}),
		},
	}
}

// Start iş akışını başlatır
func (r *WorkflowRuntime) Start(ctx context.Context) error {
	r.mutex.Lock()
	if r.state.Status != StatusPending {
		r.mutex.Unlock()
		return fmt.Errorf("iş akışı zaten başlatılmış")
	}

	r.state.Status = StatusRunning
	r.state.StartedAt = time.Now()
	r.mutex.Unlock()

	// İlk adımı başlat
	if len(r.definition.Steps) > 0 {
		r.state.CurrentStepID = r.definition.Steps[0].ID
		return r.executeCurrentStep(ctx)
	}

	return fmt.Errorf("iş akışında hiç adım yok")
}

// executeCurrentStep mevcut adımı çalıştırır
func (r *WorkflowRuntime) executeCurrentStep(ctx context.Context) error {
	r.mutex.RLock()
	currentStepID := r.state.CurrentStepID
	r.mutex.RUnlock()

	var currentStep *StepDefinition
	for _, step := range r.definition.Steps {
		if step.ID == currentStepID {
			currentStep = &step
			break
		}
	}

	if currentStep == nil {
		return fmt.Errorf("adım bulunamadı: %s", currentStepID)
	}

	// Adım için context hazırla
	stepCtx := ctx
	if currentStep.Timeout > 0 {
		var cancel context.CancelFunc
		stepCtx, cancel = context.WithTimeout(ctx, currentStep.Timeout)
		defer cancel()
	}

	// Adımı çalıştır
	result, err := r.engine.ExecuteStep(stepCtx, currentStepID, r.state.Context)

	r.mutex.Lock()

	if err != nil {
		if currentStep.RetryPolicy != nil {
			// TODO: Retry logic implementation
			r.mutex.Unlock()
			return err
		}
		r.state.Status = StatusFailed
		r.state.Error = err
		r.mutex.Unlock()
		return err
	}

	// Sonucu kaydet
	r.state.StepResults[currentStepID] = result

	// Sonraki adımı belirle
	if len(currentStep.NextSteps) > 0 {
		nextStepID := currentStep.NextSteps[0]
		r.state.CurrentStepID = nextStepID
		r.mutex.Unlock()
		return r.executeCurrentStep(ctx)
	}

	// İş akışı tamamlandı
	now := time.Now()
	r.state.CompletedAt = &now
	r.state.Status = StatusCompleted
	r.mutex.Unlock()
	return nil
}

// GetState iş akışının mevcut durumunu döndürür
func (r *WorkflowRuntime) GetState() WorkflowState {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return *r.state
}

// Cancel iş akışını iptal eder
func (r *WorkflowRuntime) Cancel() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state.Status != StatusRunning {
		return fmt.Errorf("iş akışı çalışır durumda değil")
	}

	r.state.Status = StatusCanceled
	now := time.Now()
	r.state.CompletedAt = &now
	return nil
}
