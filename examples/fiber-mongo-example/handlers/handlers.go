package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/parevo-lab/maestro"
	"github.com/parevo-lab/maestro/examples/fiber-mongo-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func SetDB(database *mongo.Database) {
	db = database
}

func CreateWorkflow(c *fiber.Ctx) error {
	workflow := new(models.Workflow)
	if err := c.BodyParser(workflow); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	workflow.ID = primitive.NewObjectID()
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()
	workflow.Status = "created"

	// Maestro engine ile workflow başlat
	eng := maestro.NewEngine()
	eng.RegisterStep(workflow.ID.Hex(), func(ctx context.Context, data interface{}) (interface{}, error) {
		return workflow, nil
	})

	_, err := db.Collection("workflows").InsertOne(context.Background(), workflow)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(workflow)
}

func GetWorkflows(c *fiber.Ctx) error {
	var workflows []models.Workflow
	cursor, err := db.Collection("workflows").Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &workflows); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(workflows)
}

func GetWorkflow(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	workflow := new(models.Workflow)
	err = db.Collection("workflows").FindOne(context.Background(), bson.M{"_id": id}).Decode(workflow)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Workflow not found"})
	}

	return c.JSON(workflow)
}

func UpdateWorkflowStatus(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	type UpdateStatus struct {
		Status string `json:"status"`
	}

	var update UpdateStatus
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := db.Collection("workflows").UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":     update.Status,
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if result.ModifiedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Workflow not found"})
	}

	return c.JSON(fiber.Map{"message": "Status updated successfully"})
}
