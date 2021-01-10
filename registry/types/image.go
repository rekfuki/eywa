package types

import "time"

// GetImagesResponse represents GET ALL response
type GetImagesResponse struct {
	Objects []Image `json:"objects"`
	Total   int     `json:"total"`
	Page    int     `json:"page_number"`
	PerPage int     `json:"per_page"`
}

// Image represents an image
type Image struct {
	ID             string    `bson:"_id" json:"id"`
	UserID         string    `bson:"user_id" json:"-"`
	TaggedRegistry string    `bson:"registry" json:"tagged_registry,omitempty"`
	Language       string    `bson:"languge" json:"language"`
	Name           string    `bson:"name" json:"name"`
	Version        string    `bson:"version" json:"version"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	State          string    `bson:"state" json:"state"`
	Source         string    `bson:"source" json:"-"`
}
