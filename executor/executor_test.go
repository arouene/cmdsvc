package executor

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteWithLogger(t *testing.T) {
	const testString = "test"

	// mock io.Writer
	buffer := &bytes.Buffer{}

	executor := WithLogger(buffer, buffer)
	err := executor.Execute("test", "echo '" + testString + "'")
	if err != nil {
		t.Errorf("Error while executing command")
	}

	// wait for the job to complete
	for {
		jobs := executor.GetJob("test")
		if len(jobs) == 0 {
			break
		}
	}

	result := buffer.String()
	if strings.TrimSpace(result) != testString {
		t.Errorf("Want %s, got %s", testString, result)
	}
}
