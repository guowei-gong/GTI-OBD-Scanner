package GTI_OBD_Scanner

import "testing"

var logger = NewLogger(
	WithStackLevel(WarnLevel),
	WithFormat(TextFormat),
	WithCallerSkip(1),
	WithClassifiedStorage(true),
)

func TestNewLogger(t *testing.T) {
	// logger.Print(InfoLevel, "hello GTI")

	logger.Info("Hello GTI")

	logger.Warn("Hello GTI")
}

func TestLog(t *testing.T) {
	Info("Hello GTI")
}
