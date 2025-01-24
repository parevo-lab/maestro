package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/parevo-lab/maestro"
	"github.com/parevo-lab/maestro/examples/fiber-mongo-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	col    *mongo.Collection
	engine *maestro.WorkflowEngine
}

func NewHandler(col *mongo.Collection, engine *maestro.WorkflowEngine) *Handler {
	return &Handler{col: col, engine: engine}
}

func (h *Handler) CreateWorkflow(c *fiber.Ctx) error {
	workflow := new(models.Workflow)
	if err := c.BodyParser(workflow); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	workflow.ID = primitive.NewObjectID()
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()
	workflow.Status = "pending"
	workflow.Results = make([]models.Result, 0)

	// Her adıma ID ve sıra numarası atama
	for i := range workflow.Steps {
		workflow.Steps[i].ID = primitive.NewObjectID()
		workflow.Steps[i].Order = i + 1
		workflow.Steps[i].Status = "pending"
		workflow.Steps[i].CreatedAt = time.Now()
		workflow.Steps[i].UpdatedAt = time.Now()
	}

	if len(workflow.Steps) > 0 {
		workflow.CurrentStep = workflow.Steps[0].ID
	}

	// Maestro workflow'unu başlat
	go func() {
		ctx := context.Background()
		for _, step := range workflow.Steps {
			// Dinamik step çalıştırma
			result, err := h.executeStep(ctx, step)
			if err != nil {
				h.updateWorkflowStatus(workflow.ID, "failed", step.ID, err.Error())
				return
			}
			h.updateWorkflowStatus(workflow.ID, "completed", step.ID, result)
		}
	}()

	_, err := h.col.InsertOne(context.Background(), workflow)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create workflow"})
	}

	return c.Status(201).JSON(workflow)
}

func (h *Handler) executeStep(ctx context.Context, step models.Step) (interface{}, error) {
	// Step tipine göre parametreleri işle
	stepResult := map[string]interface{}{
		"step_type":  step.Type,
		"parameters": step.Parameters,
		"started_at": time.Now(),
	}

	// Step'in kendi mantığını çalıştır
	switch step.Type {
	case "validation", "check", "control":
		// Kontrol adımları için
		stepResult["validation_result"] = true
		stepResult["checks_passed"] = len(step.Parameters)
	case "process", "execute", "run":
		// İşlem adımları için
		stepResult["process_status"] = "completed"
		stepResult["processed_items"] = len(step.Parameters)
	case "notification", "alert":
		// Bildirim adımları için
		stepResult["notification_sent"] = true
		stepResult["recipients"] = step.Parameters["recipients"]
	}

	stepResult["completed_at"] = time.Now()
	stepResult["duration"] = time.Since(stepResult["started_at"].(time.Time))

	return stepResult, nil
}

func (h *Handler) updateWorkflowStatus(workflowID primitive.ObjectID, status string, stepID primitive.ObjectID, result interface{}) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"results": models.Result{
				StepID:      stepID,
				Status:      status,
				Output:      result,
				CompletedAt: time.Now(),
			},
		},
	}

	_, err := h.col.UpdateOne(
		context.Background(),
		bson.M{"_id": workflowID},
		update,
	)

	return err
}

func (h *Handler) GetWorkflow(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var workflow models.Workflow
	err = h.col.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&workflow)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Workflow not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get workflow"})
	}

	return c.JSON(workflow)
}

func (h *Handler) UpdateWorkflowStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	type UpdateRequest struct {
		Status string        `json:"status"`
		StepID string        `json:"step_id"`
		Result models.Result `json:"result"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Body parse error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	stepID, err := primitive.ObjectIDFromHex(req.StepID)
	if err != nil {
		log.Printf("Step ID parse error: %v\n", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid step ID format", "details": err.Error()})
	}

	// Önce workflow'u kontrol et
	var workflow models.Workflow
	err = h.col.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&workflow)
	if err != nil {
		log.Printf("Workflow find error: %v\n", err)
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Workflow not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get workflow", "details": err.Error()})
	}

	// Step'i kontrol et
	stepFound := false
	currentStepIndex := -1
	for i, step := range workflow.Steps {
		if step.ID == stepID {
			stepFound = true
			currentStepIndex = i
			break
		}
	}

	if !stepFound {
		return c.Status(400).JSON(fiber.Map{"error": "Step not found in workflow"})
	}

	// Results alanını kontrol et ve gerekirse başlat
	if workflow.Results == nil {
		// Önce results array'ini başlat
		_, err := h.col.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$set": bson.M{"results": []models.Result{}}},
		)
		if err != nil {
			log.Printf("Results array initialization error: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to initialize results array"})
		}
	}

	// Result'ı kontrol et ve ekle
	if req.Result.StepID.IsZero() {
		req.Result.StepID = stepID
	}
	if req.Result.CompletedAt.IsZero() {
		req.Result.CompletedAt = time.Now()
	}

	update := bson.M{
		"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"results": req.Result,
		},
	}

	if req.Status == "completed" && currentStepIndex != -1 && currentStepIndex < len(workflow.Steps)-1 {
		nextStep := workflow.Steps[currentStepIndex+1]
		update["$set"].(bson.M)["current_step"] = nextStep.ID
		update["$set"].(bson.M)[fmt.Sprintf("steps.%d.status", currentStepIndex)] = "completed"
		if currentStepIndex+1 < len(workflow.Steps) {
			update["$set"].(bson.M)[fmt.Sprintf("steps.%d.status", currentStepIndex+1)] = "in_progress"
		}
	}

	result, err := h.col.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	if err != nil {
		log.Printf("MongoDB update error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update workflow",
			"details": err.Error(),
			"update":  update,
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Workflow not found"})
	}

	// Güncellenmiş workflow'u getir
	err = h.col.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&workflow)
	if err != nil {
		log.Printf("Error fetching updated workflow: %v\n", err)
		return c.Status(200).JSON(fiber.Map{"message": "Workflow updated successfully but couldn't fetch the latest state"})
	}

	return c.JSON(workflow)
}

func (h *Handler) ListWorkflows(c *fiber.Ctx) error {
	cursor, err := h.col.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to list workflows"})
	}
	defer cursor.Close(context.Background())

	var workflows []models.Workflow
	if err := cursor.All(context.Background(), &workflows); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode workflows"})
	}

	return c.JSON(workflows)
}
