package log

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	zapWrapper struct {
		logger *zap.SugaredLogger
	}
)

func NewZap() Logger {
	return &zapWrapper{
		logger: newZap(),
	}
}

func newZap() *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "_m",
		NameKey:     "logger",
		LevelKey:    "_l",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		TimeKey:     "_t",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zap.DebugLevel)

	return zap.New(core).Sugar()
}

func (z *zapWrapper) Info(msg string, args ...interface{}) {
	z.logger.Infow(msg, args...)
}

func (z *zapWrapper) Warn(msg string, args ...interface{}) {
	z.logger.Warnw(msg, args...)
}

func (z *zapWrapper) Error(msg string, args ...interface{}) {
	z.logger.Errorw(msg, args...)
}

func (z *zapWrapper) With(args ...interface{}) Logger {
	return &zapWrapper{
		logger: z.logger.With(args...),
	}
}

func (z *zapWrapper) Close() error {
	err := z.logger.Sync()
	if err == nil {
		return nil
	}

	// https://github.com/uber-go/zap/issues/328
	// очень сильно негодую по этому поводу, пока что автор не фиксит эту проблему.
	if isSyncInvalidError(err) {
		return nil
	}

	return errors.Wrap(err, "sync zap logger")
}

func isSyncInvalidError(err error) bool {
	if pathErr, ok := err.(*os.PathError); ok {
		switch pathErr.Err {
		case syscall.ENOTTY, syscall.EINVAL:
			return true
		}
	}

	return false
}
