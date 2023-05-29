package logs

type LoggerType string

type Logger interface {
	Write([]byte) (int, error)
	Close()
}

type LoggerCreator func(...interface{}) Logger

var (
	loggers map[LoggerType]LoggerCreator
)

func RegisterLogger(t LoggerType, f LoggerCreator) {
	if loggers == nil {
		loggers = make(map[LoggerType]LoggerCreator)
	}
	loggers[t] = f
}

func NewLogger(t LoggerType, args ...interface{}) (Logger, bool) {
	if f, ok := loggers[t]; ok {
		return f(args...), true
	}

	return nil, false
}
