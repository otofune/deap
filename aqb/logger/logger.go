package logger

// Logger represent aqb library logger.
// Library accepts customize logger by giving context.
type Logger interface {
	// Debugf writes debug messags
	Debugf(format string, a ...interface{}) (n int, err error)
	// WithServiceName overrides service name stored in Logger
	WithServiceName(service string) Logger
}

// NoOpLogger makes all input to waste
type NoOpLogger struct{}

// Debugf nothing output
func (NoOpLogger) Debugf(_ string, _ ...interface{}) (int, error) {
	return 0, nil
}

// WithServiceName nothing changes
func (s *NoOpLogger) WithServiceName(_ string) Logger {
	return s
}
