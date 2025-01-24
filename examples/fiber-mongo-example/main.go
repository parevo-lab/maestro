package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/parevo-lab/maestro/examples/fiber-mongo-example/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var db *mongo.Database

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client
	db = client.Database(os.Getenv("DB_NAME"))
	handlers.SetDB(db)
	log.Println("Connected to MongoDB!")
}

func main() {
	app := fiber.New()

	// Ana endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Maestro Fiber MongoDB Example",
		})
	})

	// Workflow routes
	api := app.Group("/api")
	workflows := api.Group("/workflows")

	workflows.Post("/", handlers.CreateWorkflow)
	workflows.Get("/", handlers.GetWorkflows)
	workflows.Get("/:id", handlers.GetWorkflow)
	workflows.Patch("/:id/status", handlers.UpdateWorkflowStatus)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"mongodb": "connected",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
