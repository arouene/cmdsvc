package logs

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"
)

func TestLogs(t *testing.T) {
	const loggerType = "File"
	const pathFile = "/tmp/testLogger.txt"

	token := make([]byte, 10)
	rand.Read(token)

	if fileLogger, ok := NewLogger(loggerType, pathFile); ok {
		defer fileLogger.Close()

		if _, err := fileLogger.Write(token); err != nil {
			t.Errorf("Error while writing file: %s", err)
		}

		content, err := os.ReadFile(pathFile)
		if err != nil {
			t.Errorf("Got error while reading file: %s", err)
		}

		if !bytes.Equal(content, token) {
			t.Errorf("Want %s, got %s", token, content)
		}
	} else {
		t.Errorf("Got nil Logger")
	}
}
