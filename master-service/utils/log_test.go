package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushLogf(t *testing.T) {
	logFile := "service.json"

	// Clean up created file after test
	defer os.Remove(logFile)

	// Pre-clean in case a previous failed test left it behind
	os.Remove(logFile)

	// Call the function to be tested
	PushLogf("UTC", "Test log message", "Test error message")

	// Check if the file was created
	_, err := os.Stat(logFile)
	assert.NoError(t, err, "Log file should be created")

	// Read the file content
	data, err := os.ReadFile(logFile)
	assert.NoError(t, err, "Should be able to read log file")

	// Check if the content is valid JSON and contains the log
	var logs Logs
	err = json.Unmarshal(data, &logs)
	assert.NoError(t, err, "Log file content should be valid JSON")
	assert.NotEmpty(t, logs.Logs, "Logs array should not be empty")
	assert.Contains(t, logs.Logs[0].Log, "Test log message")
	assert.Contains(t, logs.Logs[0].Log, "Test error message")
}
