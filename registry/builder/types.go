package builder

// States of image/build
const (
	StateBuilding = "building"
	StateSuccess  = "success"
	StateFailed   = "failed"
	StateQueued   = "queued"
)

// BuildRequest represents the payload of a request to build a new image
type BuildRequest struct {
	ImageID        string
	UserID         string
	Name           string
	Runtime        string
	Version        string
	ZippedSource   []byte
	tmpDir         string
	LogFile        string
	ExecutablePath *string
	requiredFiles  []string
}
