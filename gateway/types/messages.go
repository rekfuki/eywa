package types

import (
	"fmt"
	"time"
)

// SyncExecutionStartMessage returns sync execution start message to be logged
func SyncExecutionStartMessage(functionName string) string {
	return fmt.Sprintf("SYNC EXECUTION STARTED: Function %q execution started", functionName)
}

// SyncExecutionFinishMessage returns sync execution finish message to be logged
func SyncExecutionFinishMessage(functionName string, status int, duration time.Duration) string {
	return fmt.Sprintf("SYNC EXECUTION FINISHED: Function %q execution took %d ms, status %d", functionName, duration.Milliseconds(), status)
}

// AsyncExecutionStartMessage returns async execution start message to be logged
func AsyncExecutionStartMessage(attempt int, functionName string) string {
	return fmt.Sprintf("ASYNC ATTEMPT #%d STARTED: Function %q execution started", attempt, functionName)
}

// AsyncExecutionFinishMessage returns async execution end message to be logged
func AsyncExecutionFinishMessage(functionName string, attempt, status int, duration time.Duration) string {
	return fmt.Sprintf("ASYNC ATTEMPT #%d FINISHED: Function %q execution took %d ms, status: %d", attempt, functionName, duration.Milliseconds(), status)
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

// CallbackError returns callback errror message to be logged
func CallbackError(message string) string {
	return fmt.Sprintf("CALLBACK ERROR: %s", message)
}

// QueuedMessage returns queued message to be logged
func QueuedMessage(requestID, functionID, functionName string) string {
	return fmt.Sprintf("QUEUED: Request ID: %q Function ID: %q Function Name: %q",
		requestID, functionID, functionName)
}
