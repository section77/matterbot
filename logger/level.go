package logger

// LogLevel defines all currently supported logging levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
	Disabled
)

var logLevel = InfoLevel

func (level LogLevel) String() string {
	var str string
	switch level {
	case DebugLevel:
		str = "DEBUG"
	case InfoLevel:
		str = "INFO"
	case ErrorLevel:
		str = "ERROR"
	}
	return str
}

func SetLogLevel(level LogLevel) {
	mutex.Lock()
	logLevel = level
	mutex.Unlock()
}
