package engine

import (
	"context"
	"testing"
	"time"
)

func TestNewWorkflowRuntime(t *testing.T) {
	engine := NewWorkflowEngine()
	definition := NewWorkflowDefinition("test", "Test Workflow", "Test Description")

	runtime := NewWorkflowRuntime(engine, definition)
	if runtime == nil {
		t.Error("NewWorkflowRuntime should return a non-nil runtime")
	}

	state := runtime.GetState()
	if state.Status != StatusPending {
		t.Error("Initial state should be pending")
	}
}

func TestWorkflowExecution(t *testing.T) {
	engine := NewWorkflowEngine()
	definition := NewWorkflowDefinition("test", "Test Workflow", "Test Description")

	// Test adımlarını oluştur
	step1Results := make(chan interface{}, 1)
	step2Results := make(chan interface{}, 1)

	engine.RegisterStep("step1", func(ctx context.Context, data interface{}) (interface{}, error) {
		result := "step1-complete"
		step1Results <- result
		return result, nil
	})

	engine.RegisterStep("step2", func(ctx context.Context, data interface{}) (interface{}, error) {
		result := "step2-complete"
		step2Results <- result
		return result, nil
	})

	// İş akışını tanımla
	step1 := NewStepDefinition("step1", "First Step", StepTypeTask).
		WithNextSteps("step2")
	step2 := NewStepDefinition("step2", "Second Step", StepTypeTask)

	definition.AddStep(step1)
	definition.AddStep(step2)

	// Runtime oluştur ve başlat
	runtime := NewWorkflowRuntime(engine, definition)

	// İş akışını başlat
	err := runtime.Start(context.Background())
	if err != nil {
		t.Fatalf("Workflow execution failed: %v", err)
	}

	// Adım 1'in tamamlanmasını bekle
	select {
	case result := <-step1Results:
		if result != "step1-complete" {
			t.Error("Step 1 returned incorrect result")
		}
	case <-time.After(time.Second):
		t.Error("Step 1 execution timeout")
	}

	// Adım 2'nin tamamlanmasını bekle
	select {
	case result := <-step2Results:
		if result != "step2-complete" {
			t.Error("Step 2 returned incorrect result")
		}
	case <-time.After(time.Second):
		t.Error("Step 2 execution timeout")
	}

	// Son durumu kontrol et
	finalState := runtime.GetState()
	if finalState.Status != StatusCompleted {
		t.Error("Workflow should be completed")
	}

	if finalState.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
}

func TestWorkflowCancellation(t *testing.T) {
	engine := NewWorkflowEngine()
	definition := NewWorkflowDefinition("test", "Test Workflow", "Test Description")

	// Uzun süren bir adım oluştur
	engine.RegisterStep("long-step", func(ctx context.Context, data interface{}) (interface{}, error) {
		time.Sleep(2 * time.Second)
		return "complete", nil
	})

	step := NewStepDefinition("long-step", "Long Step", StepTypeTask)
	definition.AddStep(step)

	runtime := NewWorkflowRuntime(engine, definition)

	// İş akışını başlat
	go func() {
		err := runtime.Start(context.Background())
		if err != nil {
			t.Errorf("Workflow execution failed: %v", err)
		}
	}()

	// Kısa bir süre bekle ve iptal et
	time.Sleep(100 * time.Millisecond)
	err := runtime.Cancel()
	if err != nil {
		t.Errorf("Workflow cancellation failed: %v", err)
	}

	// Son durumu kontrol et
	finalState := runtime.GetState()
	if finalState.Status != StatusCanceled {
		t.Error("Workflow should be canceled")
	}
}

func TestWorkflowTimeout(t *testing.T) {
	engine := NewWorkflowEngine()
	definition := NewWorkflowDefinition("test", "Test Workflow", "Test Description")

	// Timeout'lu bir adım oluştur
	engine.RegisterStep("timeout-step", func(ctx context.Context, data interface{}) (interface{}, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			return "complete", nil
		}
	})

	step := NewStepDefinition("timeout-step", "Timeout Step", StepTypeTask).
		WithTimeout(100 * time.Millisecond)
	definition.AddStep(step)

	runtime := NewWorkflowRuntime(engine, definition)

	// İş akışını başlat
	err := runtime.Start(context.Background())
	if err == nil {
		t.Error("Workflow should fail with timeout")
	}

	// Son durumu kontrol et
	finalState := runtime.GetState()
	if finalState.Status != StatusFailed {
		t.Error("Workflow should be failed")
	}
}
