package main

import (
	"context"
	"fmt"
	"time"

	"github.com/parevo-lab/maestro"
)

// FileShare represents a file sharing request
type FileShare struct {
	Files    []File
	Users    []User
	Approved bool
}

type File struct {
	ID       string
	Name     string
	Size     int64
	MimeType string
}

type User struct {
	ID    string
	Email string
	Name  string
}

func main() {
	// Create a new workflow for file sharing
	wfEngine := maestro.NewEngine()

	// Add observer for error handling
	wfEngine.AddObserver(func(event maestro.Event) {
		if event.Type == maestro.EventStepFailed {
			fmt.Printf("Workflow error at step %s: %v\n", event.StepID, event.Data)
		}
	})

	// Register steps
	wfEngine.RegisterStep("fetch-files", func(ctx context.Context, data interface{}) (interface{}, error) {
		// Simulate fetching files from storage
		files := []File{
			{
				ID:       "file1",
				Name:     "document.pdf",
				Size:     1024 * 1024,
				MimeType: "application/pdf",
			},
			{
				ID:       "file2",
				Name:     "image.jpg",
				Size:     2048 * 1024,
				MimeType: "image/jpeg",
			},
		}

		return &FileShare{Files: files}, nil
	})

	// Step 2: Fetch users
	wfEngine.RegisterStep("fetch-users", func(ctx context.Context, data interface{}) (interface{}, error) {
		fileShare := data.(*FileShare)

		// Simulate fetching users from database
		users := []User{
			{
				ID:    "user1",
				Email: "john@example.com",
				Name:  "John Doe",
			},
			{
				ID:    "user2",
				Email: "jane@example.com",
				Name:  "Jane Smith",
			},
		}

		fileShare.Users = users
		return fileShare, nil
	})

	// Step 3: Send notifications
	wfEngine.RegisterStep("send-notifications", func(ctx context.Context, data interface{}) (interface{}, error) {
		fileShare := data.(*FileShare)

		// Simulate sending notifications
		fmt.Println("Sending notifications to users:")
		for _, user := range fileShare.Users {
			fmt.Printf("- Sending notification to %s (%s)\n", user.Name, user.Email)
			fmt.Printf("  Files to be shared:\n")
			for _, file := range fileShare.Files {
				fmt.Printf("  - %s (%s)\n", file.Name, file.MimeType)
			}
		}

		return fileShare, nil
	})

	// Step 4: Wait for approval
	wfEngine.RegisterStep("wait-for-approval", func(ctx context.Context, data interface{}) (interface{}, error) {
		fileShare := data.(*FileShare)

		// Simulate waiting for approval with a timeout
		approvalChan := make(chan bool)
		timeoutChan := time.After(24 * time.Hour)

		// Simulate approval process in background
		go func() {
			// In a real application, this would wait for user input or external system
			time.Sleep(5 * time.Second)
			approvalChan <- true
		}()

		// Wait for either approval or timeout
		select {
		case approved := <-approvalChan:
			fileShare.Approved = approved
			if approved {
				fmt.Println("File share request approved!")
				fmt.Println("Granting access to files...")
				for _, user := range fileShare.Users {
					fmt.Printf("- Granted access to %s\n", user.Name)
				}
			} else {
				fmt.Println("File share request rejected.")
			}
		case <-timeoutChan:
			return nil, fmt.Errorf("approval timeout exceeded")
		}

		return fileShare, nil
	})

	// Execute steps in sequence
	var result interface{}
	var err error

	result, err = wfEngine.ExecuteStep(context.Background(), "fetch-files", nil)
	if err != nil {
		fmt.Printf("Workflow failed at fetch-files: %v\n", err)
		return
	}

	result, err = wfEngine.ExecuteStep(context.Background(), "fetch-users", result)
	if err != nil {
		fmt.Printf("Workflow failed at fetch-users: %v\n", err)
		return
	}

	result, err = wfEngine.ExecuteStep(context.Background(), "send-notifications", result)
	if err != nil {
		fmt.Printf("Workflow failed at send-notifications: %v\n", err)
		return
	}

	result, err = wfEngine.ExecuteStep(context.Background(), "wait-for-approval", result)
	if err != nil {
		fmt.Printf("Workflow failed at wait-for-approval: %v\n", err)
		return
	}

	// Print final result
	fileShare := result.(*FileShare)
	fmt.Printf("\nWorkflow completed!\n")
	fmt.Printf("Files shared: %d\n", len(fileShare.Files))
	fmt.Printf("Users notified: %d\n", len(fileShare.Users))
	fmt.Printf("Approved: %v\n", fileShare.Approved)
}
