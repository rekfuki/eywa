package builder

import "fmt"

// BuildQueuedMessage return build queued message
func BuildQueuedMessage(id, language, version string) string {
	return fmt.Sprintf(
		"########## BUILD QUEUED ##########\n"+
			"VERSION: %s\n"+
			"LANGUAGE: %s\n"+
			"ID: %s\n", version, language, id,
	)
}

// BuildStartMessage return build start message
func BuildStartMessage() string {
	return "########## BUILD START ##########\n"
}

// BuildSystemErrorMessage returns system error message
func BuildSystemErrorMessage(message string) string {
	return fmt.Sprintf("SYSTEM ERROR: %s\n", message)
}

// BuildUserErrorMessage returns user error message
func BuildUserErrorMessage(message string) string {
	return fmt.Sprintf("USER ERROR: %s\n", message)
}

// BuildErrorMessage returns generic error message
func BuildErrorMessage(message string) string {
	return fmt.Sprintf("ERROR: %s\n", message)
}

// BuildSuccessMessage returns success message
func BuildSuccessMessage() string {
	return "########## BUILD FINISHED ##########"
}

// BuildFailedMessage return build failed message
func BuildFailedMessage() string {
	return "########## BUILD FAILED ##########"
}
