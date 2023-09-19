package graphqlws

type Logger interface {
	Infof(format string, data ...any)
	Debugf(format string, data ...any)
	Errorf(format string, data ...any)
	Warnf(format string, data ...any)
}

type noopLogger struct{}

func (n *noopLogger) Infof(format string, data ...any)  {}
func (n *noopLogger) Debugf(format string, data ...any) {}
func (n *noopLogger) Errorf(format string, data ...any) {}
func (n *noopLogger) Warnf(format string, data ...any)  {}
