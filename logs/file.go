package logs

import (
	"log"
	"os"
)

func init() {
	RegisterLogger("File", NewFileLogger)
}

type FileLogger struct {
	file *os.File
}

func NewFileLogger(args ...interface{}) Logger {
	log.Println("New file Logger")

	if len(args) != 1 {
		panic("NewFileLogger needs one argument (path)")
	}
	filename := args[0].(string)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return FileLogger{
		file: f,
	}
}

func (f FileLogger) Write(b []byte) (int, error) {
	return f.file.Write(b)
}

func (f FileLogger) Close() {
	f.file.Close()
}
