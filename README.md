# Maestro

[![Go Reference](https://pkg.go.dev/badge/github.com/parevo-lab/maestro.svg)](https://pkg.go.dev/github.com/parevo-lab/maestro)
[![Go Report Card](https://goreportcard.com/badge/github.com/parevo-lab/maestro)](https://goreportcard.com/report/github.com/parevo-lab/maestro)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Maestro is a powerful Go package for building and managing workflow orchestration systems. It provides a simple, flexible, and type-safe way to define, execute, and monitor workflows in Go applications.

## Features

- 🎯 Type-safe workflow definitions
- 🔄 Flexible workflow composition
- 👥 Concurrent workflow execution
- ✅ Built-in error handling and recovery
- 📊 Workflow state management
- 🔍 Progress tracking and monitoring
- 🛡️ Context-aware execution
- 🚀 High performance and low overhead
- 📦 Zero external dependencies

## Installation

```bash
go get github.com/parevo-lab/maestro
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/parevo-lab/maestro"
)

func main() {
    // Create a new workflow
    flow := maestro.NewWorkflow("data-processing")

    // Define workflow steps
    flow.AddStep("fetch-data", func(ctx context.Context, data interface{}) (interface{}, error) {
        // Fetch data implementation
        return []string{"data1", "data2"}, nil
    })

    flow.AddStep("process-data", func(ctx context.Context, data interface{}) (interface{}, error) {
        items := data.([]string)
        // Process data implementation
        return items, nil
    })

    // Execute workflow
    result, err := flow.Execute(context.Background(), nil)
    if err != nil {
        fmt.Printf("Workflow failed: %v\n", err)
        return
    }

    fmt.Printf("Workflow completed: %v\n", result)
}
```

## Core Concepts

### Workflow

A workflow is a sequence of steps that are executed in a defined order. Each workflow has:

- A unique identifier
- A collection of steps
- Input and output types
- Execution context
- Error handling mechanisms

```go
type Workflow struct {
    ID      string
    Steps   []Step
    Context context.Context
}
```

### Step

A step is a single unit of work within a workflow:

```go
type Step struct {
    ID       string
    Handler  StepHandler
    Options  StepOptions
}

type StepHandler func(ctx context.Context, input interface{}) (interface{}, error)
```

### Advanced Usage

#### Parallel Execution

```go
func main() {
    flow := maestro.NewWorkflow("parallel-processing")

    // Add parallel steps
    flow.AddParallelSteps(
        maestro.Step{
            ID: "step1",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                return "result1", nil
            },
        },
        maestro.Step{
            ID: "step2",
            Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
                return "result2", nil
            },
        },
    )

    results, err := flow.Execute(context.Background(), nil)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Results: %v\n", results)
}
```

#### Error Handling and Retries

```go
flow := maestro.NewWorkflow("error-handling")

// Global error handler
flow.OnError(func(err error) error {
    return fmt.Errorf("workflow error: %w", err)
})

// Step with retry policy
flow.AddStep("risky-operation", func(ctx context.Context, data interface{}) (interface{}, error) {
    // Implementation with error handling
    return nil, nil
}).WithRetry(&maestro.RetryPolicy{
    MaxAttempts: 3,
    Delay: time.Second,
    BackoffFactor: 2.0,
})
```

#### State Management

```go
flow := maestro.NewWorkflow("stateful-workflow")

// Add state store
flow.WithStateStore(maestro.NewInMemoryStore())

// Access state in steps
flow.AddStep("stateful-step", func(ctx context.Context, data interface{}) (interface{}, error) {
    state := maestro.GetState(ctx)
    state.Set("key", "value")
    return state.Get("key"), nil
})
```

#### Conditional Workflows

```go
flow := maestro.NewWorkflow("conditional-flow")

flow.AddStep("check-condition", func(ctx context.Context, data interface{}) (interface{}, error) {
    if someCondition {
        return flow.ExecuteBranch("success-branch")
    }
    return flow.ExecuteBranch("failure-branch")
})

flow.AddBranch("success-branch", maestro.NewWorkflow("success-flow"))
flow.AddBranch("failure-branch", maestro.NewWorkflow("failure-flow"))
```

## Examples

For more examples, check out the [examples](examples) directory in the repository.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
