package logger

// Logger is an interface for logging in long-running processes like updating
// holidays storage
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

// NilLogger is a logger doing nothing
type NilLogger struct{}

func (n *NilLogger) Debugf(format string, args ...interface{})   {}
func (n *NilLogger) Infof(format string, args ...interface{})    {}
func (n *NilLogger) Warningf(format string, args ...interface{}) {}
func (n *NilLogger) Errorf(format string, args ...interface{})   {}

func (n *NilLogger) Debug(args ...interface{})   {}
func (n *NilLogger) Info(args ...interface{})    {}
func (n *NilLogger) Warning(args ...interface{}) {}
func (n *NilLogger) Error(args ...interface{})   {}

var _ Logger = &NilLogger{}
