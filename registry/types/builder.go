package types

import "archive/zip"

const (
	// StateBuilding ...
	StateBuilding = "building"
	// StateSuccess ...
	StateSuccess = "success"
	// StateFailed ...
	StateFailed = "failed"
)

// BuildRequest represents the payload of a request to build a new image
type BuildRequest struct {
	ID        string
	Language  string
	Version   string
	ZipReader *zip.Reader
}
