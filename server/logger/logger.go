package logger

type Logger interface {
	Infof(format string, data ...any)
	Debugf(format string, data ...any)
	Errorf(format string, data ...any)
	Warnf(format string, data ...any)
}

type NoopLogger struct{}

func (n *NoopLogger) Infof(format string, data ...any)  {}
func (n *NoopLogger) Debugf(format string, data ...any) {}
func (n *NoopLogger) Errorf(format string, data ...any) {}
func (n *NoopLogger) Warnf(format string, data ...any)  {}
