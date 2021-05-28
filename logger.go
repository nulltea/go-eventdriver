package eventdriver

// Logger defines common logging interface.
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

// NopLogger is a default implementation of the Logger interface.
//
// It will do nothing.
type NopLogger struct { }

func (n NopLogger) Error(args ...interface{}) {}

func (n NopLogger) Errorf(format string, args ...interface{}) {}

func (n NopLogger) Warning(args ...interface{}) {}

func (n NopLogger) Warningf(format string, args ...interface{}) {}

var _logger Logger = NopLogger{}
