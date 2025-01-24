package engine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkflowEngine iş akışı motorunun ana yapısı
type WorkflowEngine struct {
	steps     map[string]StepFunc
	mutex     sync.RWMutex
	observers []ObserverFunc
}

// StepFunc bir iş akışı adımını temsil eden fonksiyon tipi
type StepFunc func(ctx context.Context, data interface{}) (interface{}, error)

// ObserverFunc iş akışı olaylarını dinleyen fonksiyon tipi
type ObserverFunc func(event Event)

// Event iş akışındaki olayları temsil eder
type Event struct {
	Type      EventType
	StepID    string
	Data      interface{}
	Timestamp time.Time
}

// EventType olay tiplerini temsil eder
type EventType string

const (
	EventStepStarted  EventType = "step_started"
	EventStepComplete EventType = "step_completed"
	EventStepFailed   EventType = "step_failed"
)

// NewWorkflowEngine yeni bir iş akışı motoru oluşturur
func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		steps:     make(map[string]StepFunc),
		observers: make([]ObserverFunc, 0),
	}
}

// RegisterStep yeni bir adım kaydeder
func (e *WorkflowEngine) RegisterStep(id string, step StepFunc) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.steps[id] = step
}

// AddObserver yeni bir gözlemci ekler
func (e *WorkflowEngine) AddObserver(observer ObserverFunc) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.observers = append(e.observers, observer)
}

// notifyObservers tüm gözlemcilere olay bildirir
func (e *WorkflowEngine) notifyObservers(event Event) {
	for _, observer := range e.observers {
		observer(event)
	}
}

// ExecuteStep belirli bir adımı çalıştırır
func (e *WorkflowEngine) ExecuteStep(ctx context.Context, stepID string, data interface{}) (interface{}, error) {
	e.mutex.RLock()
	step, exists := e.steps[stepID]
	e.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("adım bulunamadı: %s", stepID)
	}

	// Adım başlangıç olayını bildir
	e.notifyObservers(Event{
		Type:      EventStepStarted,
		StepID:    stepID,
		Data:      data,
		Timestamp: time.Now(),
	})

	result, err := step(ctx, data)
	if err != nil {
		// Hata olayını bildir
		e.notifyObservers(Event{
			Type:      EventStepFailed,
			StepID:    stepID,
			Data:      err,
			Timestamp: time.Now(),
		})
		return nil, err
	}

	// Başarılı tamamlanma olayını bildir
	e.notifyObservers(Event{
		Type:      EventStepComplete,
		StepID:    stepID,
		Data:      result,
		Timestamp: time.Now(),
	})

	return result, nil
}
