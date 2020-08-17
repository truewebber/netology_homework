package log

type (
	Logger interface {
		Info(msg string, args ...interface{})
		Warn(msg string, args ...interface{})
		Error(msg string, args ...interface{})

		With(args ...interface{}) Logger

		Close() error
	}
)
