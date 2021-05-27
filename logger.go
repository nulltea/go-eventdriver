package eventdriver

// Logger defines common logging interface.
type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}
