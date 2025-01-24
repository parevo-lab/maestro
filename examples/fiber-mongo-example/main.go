package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/parevo-lab/maestro"
	"github.com/parevo-lab/maestro/examples/fiber-mongo-example/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// MongoDB bağlantısı
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Maestro engine'i başlat
	engine := maestro.NewEngine()

	// Hata gözlemcisi ekle
	engine.AddObserver(func(event maestro.Event) {
		if event.Type == maestro.EventStepFailed {
			log.Printf("Workflow error at step %s: %v\n", event.StepID, event.Data)
		}
	})

	// Fiber app oluşturma
	app := fiber.New()

	// Collection
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "workflow_db"
	}

	db := client.Database(dbName)
	col := db.Collection("workflows")

	// Handler oluşturma
	h := handlers.NewHandler(col, engine)

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	workflows := v1.Group("/workflows")
	workflows.Post("/", h.CreateWorkflow)
	workflows.Get("/", h.ListWorkflows)
	workflows.Get("/:id", h.GetWorkflow)
	workflows.Put("/:id", h.UpdateWorkflowStatus)

	// Server başlatma
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
