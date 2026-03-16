package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLEntry struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	URL         string             `json:"url"`
	ShortCode   string             `json:"shortCode"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	AccessCount int                `json:"accessCount"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}
