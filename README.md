# 🎭 Maestro

[![Go Reference](https://pkg.go.dev/badge/github.com/parevo-lab/maestro.svg)](https://pkg.go.dev/github.com/parevo-lab/maestro)
[![Go Report Card](https://goreportcard.com/badge/github.com/parevo-lab/maestro)](https://goreportcard.com/report/github.com/parevo-lab/maestro)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Maestro is a lightweight, high-performance workflow orchestration engine for Go applications. It provides a robust foundation for building complex, distributed workflows while maintaining simplicity and type safety.

## ✨ Key Features

- 🔄 **Simple Workflow Definition**: Easy-to-use API for defining workflow steps and their dependencies
- 🛡️ **Type Safety**: Full Go type system support for workflow data
- 📊 **Event-Driven Architecture**: Built-in observer pattern for workflow monitoring
- 🔒 **Thread Safety**: Concurrent execution with proper synchronization
- 🎯 **Context Awareness**: Native support for Go context for timeout and cancellation
- 📈 **Extensible**: Easy to add custom steps and observers
- 🚀 **High Performance**: Minimal overhead and efficient execution
- 🧪 **Testing Support**: Designed with testability in mind

## 📦 Installation

```bash
go get github.com/parevo-lab/maestro
```

## 🚀 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/parevo-lab/maestro"
)

// Define your workflow data structures
type Task struct {
    Name        string
    Description string
    IsCompleted bool
}

func main() {
    // Create a new workflow engine
    engine := maestro.NewEngine()

    // Add error handling observer
    engine.AddObserver(func(event maestro.Event) {
        if event.Type == maestro.EventStepFailed {
            fmt.Printf("Workflow error at step %s: %v\n", event.StepID, event.Data)
        }
    })

    // Register workflow steps
    engine.RegisterStep("create-task", func(ctx context.Context, data interface{}) (interface{}, error) {
        // Create a new task
        task := &Task{
            Name:        "Sample Task",
            Description: "This is a sample task",
            IsCompleted: false,
        }
        return task, nil
    })

    engine.RegisterStep("process-task", func(ctx context.Context, data interface{}) (interface{}, error) {
        // Process the task
        task := data.(*Task)
        task.IsCompleted = true
        return task, nil
    })

    // Execute workflow steps
    ctx := context.Background()
    
    // Execute first step
    result, err := engine.ExecuteStep(ctx, "create-task", nil)
    if err != nil {
        fmt.Printf("Error creating task: %v\n", err)
        return
    }

    // Execute second step
    result, err = engine.ExecuteStep(ctx, "process-task", result)
    if err != nil {
        fmt.Printf("Error processing task: %v\n", err)
        return
    }

    // Print final result
    task := result.(*Task)
    fmt.Printf("Task completed: %s - Completed: %v\n", task.Name, task.IsCompleted)
}
```

## 📚 Core Concepts

### WorkflowEngine

The `WorkflowEngine` is the central component that manages workflow execution:

```go
type WorkflowEngine struct {
    steps     map[string]StepFunc
    observers []ObserverFunc
}
```

### Steps

Steps are the building blocks of workflows:

```go
type StepFunc func(ctx context.Context, data interface{}) (interface{}, error)
```

### Events

The engine emits events during workflow execution:

```go
type Event struct {
    Type      EventType
    StepID    string
    Data      interface{}
    Timestamp time.Time
}
```

## 🎯 Use Cases

- **Data Processing Pipelines**: Build complex data transformation workflows
- **Business Process Automation**: Automate multi-step business processes
- **Microservices Orchestration**: Coordinate multiple service calls
- **Task Scheduling**: Create dependent task execution flows
- **File Processing**: Handle multi-stage file processing workflows

## 📖 Examples

Check out the [examples](./examples) directory for more detailed examples, including:
- File sharing workflows
- Data processing pipelines
- Service orchestration patterns

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
