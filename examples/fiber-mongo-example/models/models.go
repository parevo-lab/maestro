package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Step struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Type        string                 `json:"type" bson:"type"`
	Status      string                 `json:"status" bson:"status"`
	Order       int                    `json:"order" bson:"order"`
	Parameters  map[string]interface{} `json:"parameters" bson:"parameters"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

type Result struct {
	StepID      primitive.ObjectID `json:"step_id" bson:"step_id"`
	Status      string             `json:"status" bson:"status"`
	Output      interface{}        `json:"output" bson:"output"`
	Error       string             `json:"error,omitempty" bson:"error,omitempty"`
	CompletedAt time.Time          `json:"completed_at" bson:"completed_at"`
}

type Workflow struct {
	ID          primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	Type        string                 `json:"type" bson:"type"`
	CreatedBy   primitive.ObjectID     `json:"created_by" bson:"created_by"`
	Status      string                 `json:"status" bson:"status"`
	Steps       []Step                 `json:"steps" bson:"steps"`
	CurrentStep primitive.ObjectID     `json:"current_step" bson:"current_step"`
	Results     []Result               `json:"results" bson:"results"`
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}
