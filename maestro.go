package maestro

import "github.com/parevo-lab/maestro/pkg/engine"

// Re-export engine types and functions
type WorkflowEngine = engine.WorkflowEngine
type StepFunc = engine.StepFunc
type ObserverFunc = engine.ObserverFunc
type Event = engine.Event
type EventType = engine.EventType

// Re-export event constants
const (
	EventStepStarted  = engine.EventStepStarted
	EventStepComplete = engine.EventStepComplete
	EventStepFailed   = engine.EventStepFailed
)

// NewEngine creates a new workflow engine
func NewEngine() *WorkflowEngine {
	return engine.NewWorkflowEngine()
}
