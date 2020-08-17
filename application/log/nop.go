package log

type (
	nop struct{}
)

func NewNop() Logger {
	return new(nop)
}

func (nop) Info(string, ...interface{}) {}

func (nop) Warn(string, ...interface{}) {}

func (nop) Error(string, ...interface{}) {}

func (nop) With(...interface{}) Logger {
	return new(nop)
}

func (nop) Close() error { return nil }
