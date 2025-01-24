package engine

import (
	"context"
	"testing"
	"time"
)

func TestNewWorkflowEngine(t *testing.T) {
	engine := NewWorkflowEngine()
	if engine == nil {
		t.Error("NewWorkflowEngine should return a non-nil engine")
	}
	if len(engine.steps) != 0 {
		t.Error("New engine should have no steps")
	}
	if len(engine.observers) != 0 {
		t.Error("New engine should have no observers")
	}
}

func TestRegisterStep(t *testing.T) {
	engine := NewWorkflowEngine()
	stepID := "test-step"

	// Test step kaydı
	engine.RegisterStep(stepID, func(ctx context.Context, data interface{}) (interface{}, error) {
		return "success", nil
	})

	if len(engine.steps) != 1 {
		t.Error("Engine should have exactly one step")
	}

	if _, exists := engine.steps[stepID]; !exists {
		t.Error("Registered step should exist in engine")
	}
}

func TestAddObserver(t *testing.T) {
	engine := NewWorkflowEngine()
	eventChan := make(chan Event, 1)

	// Test observer kaydı
	engine.AddObserver(func(event Event) {
		eventChan <- event
	})

	if len(engine.observers) != 1 {
		t.Error("Engine should have exactly one observer")
	}

	// Test event bildirimi
	testEvent := Event{
		Type:      EventStepStarted,
		StepID:    "test",
		Data:      "test-data",
		Timestamp: time.Now(),
	}

	engine.notifyObservers(testEvent)

	select {
	case receivedEvent := <-eventChan:
		if receivedEvent.Type != testEvent.Type {
			t.Error("Received event type does not match")
		}
		if receivedEvent.StepID != testEvent.StepID {
			t.Error("Received event stepID does not match")
		}
	case <-time.After(time.Second):
		t.Error("Observer notification timeout")
	}
}

func TestExecuteStep(t *testing.T) {
	engine := NewWorkflowEngine()
	stepID := "test-step"
	testData := "input-data"
	expectedResult := "success"

	// Test step kaydı
	engine.RegisterStep(stepID, func(ctx context.Context, data interface{}) (interface{}, error) {
		if data != testData {
			t.Error("Step received incorrect input data")
		}
		return expectedResult, nil
	})

	// Test step çalıştırma
	result, err := engine.ExecuteStep(context.Background(), stepID, testData)
	if err != nil {
		t.Errorf("ExecuteStep returned unexpected error: %v", err)
	}

	if result != expectedResult {
		t.Errorf("ExecuteStep returned incorrect result. Expected %v, got %v", expectedResult, result)
	}

	// Test var olmayan step
	_, err = engine.ExecuteStep(context.Background(), "non-existent", nil)
	if err == nil {
		t.Error("ExecuteStep should return error for non-existent step")
	}
}

func TestConcurrentStepExecution(t *testing.T) {
	engine := NewWorkflowEngine()
	stepID := "concurrent-step"
	executionCount := 0

	engine.RegisterStep(stepID, func(ctx context.Context, data interface{}) (interface{}, error) {
		executionCount++
		time.Sleep(10 * time.Millisecond)
		return executionCount, nil
	})

	// Concurrent execution
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := engine.ExecuteStep(context.Background(), stepID, nil)
			if err != nil {
				t.Error("Concurrent execution failed")
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	if executionCount != numGoroutines {
		t.Errorf("Expected %d executions, got %d", numGoroutines, executionCount)
	}
}
