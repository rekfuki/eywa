package types

import "time"

// GetImagesResponse represents GET ALL response
type GetImagesResponse struct {
	Objects []Image `json:"objects"`
	Total   int     `json:"total_count"`
	Page    int     `json:"page_number"`
	PerPage int     `json:"per_page"`
}

// ImageBuildResponse is returned when build request is issued
type ImageBuildResponse struct {
	BuildID   string    `json:"build_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Image represents an image
type Image struct {
	ID             string    `db:"id" json:"id"`
	UserID         string    `db:"user_id" json:"-"`
	TaggedRegistry string    `db:"registry" json:"tagged_registry,omitempty"`
	Runtime        string    `db:"language" json:"runtime"`
	Name           string    `db:"name" json:"name"`
	Version        string    `db:"version" json:"version"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	State          string    `db:"state" json:"state"`
	Size           int       `db:"size" json:"size"`
	Source         string    `db:"-" json:"-"`
}

// ImageLogs represents image build logs responsej
type ImageLogs struct {
	Logs []string `json:"logs"`
}
