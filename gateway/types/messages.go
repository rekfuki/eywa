package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StartMessage returns start message to be logged
func StartMessage(requestID, functionID, functionName string) string {
	return fmt.Sprintf("START: Request ID: %q Function ID: %q Function Name: %q",
		requestID, functionID, functionName)
}

// AttemptMessage returns attempt message to be logged
func AttemptMessage(attempt int, requestID, functionID, functionName string) string {
	return fmt.Sprintf("ATTEMPT #%d: Request ID: %q Function ID: %q Function Name: %q",
		attempt, requestID, functionID, functionName)
}

// RequestContextMessage returns request context message to be logged
func RequestContextMessage(url, query string, payload []byte, headers http.Header) string {
	bytes, err := buildContext(url, query, payload, headers, 0)
	message := string(bytes)
	if err != nil {
		message = "Failed to build context"
	}
	return fmt.Sprintf("REQUEST CONTEXT: %s", message)
}

// ServerErrorMessage returns server error message to be logged
func ServerErrorMessage() string {
	return "ERROR: System error. Please try again."
}

// FunctionNotFoundMessage returns function not found message to be logged
func FunctionNotFoundMessage(requestID, functionID, functionName string) string {
	return fmt.Sprintf("ERROR: Function not found. RequestID %q Function ID: %q Function Name: %q",
		requestID, functionID, functionName)
}

// StdoutMessage returns stdouout message to be logged
func StdoutMessage(message string) string {
	return fmt.Sprintf("STDOUT: %s", message)
}

// StderrMessage returns stderr message to be logged
func StderrMessage(message string) string {
	return fmt.Sprintf("STDERR: %s", message)
}

// CallbackError returns callback errror message to be logged
func CallbackError(message string) string {
	return fmt.Sprintf("CALLBACK ERROR: %s", message)
}

// ResponseContextMessage returns response context message to be logged
func ResponseContextMessage(url, query string, payload []byte, headers http.Header, status int) string {
	bytes, err := buildContext(url, query, payload, headers, status)
	message := string(bytes)
	if err != nil {
		message = "Failed to build context"
	}
	return fmt.Sprintf("RESPONSE CONTEXT: %s", message)
}

// QueuedMessage returns queued message to be logged
func QueuedMessage(requestID, functionID, functionName string) string {
	return fmt.Sprintf("QUEUED: Request ID: %q Function ID: %q Function Name: %q",
		requestID, functionID, functionName)
}

func buildContext(url, query string, payload []byte, headers http.Header, status int) ([]byte, error) {
	jsonContext := struct {
		URL     string      `json:"url"`
		Query   string      `json:"query"`
		Payload string      `json:"payload"`
		Headers http.Header `json:"headers"`
		Status  int         `json:"status,omitempty"`
	}{
		Query:   query,
		URL:     url,
		Payload: string(payload),
		Headers: headers,
		Status:  status,
	}

	return json.Marshal(jsonContext)
}
