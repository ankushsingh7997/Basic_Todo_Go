package todoMod

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoModel struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty"`
	Completed bool               `json:"completed,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
}
