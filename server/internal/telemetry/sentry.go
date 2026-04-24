package telemetry

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

const flushTimeout = 2 * time.Second

func Init() bool {
	dsn := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if dsn == "" {
		return false
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      firstNonEmpty(os.Getenv("APP_ENV"), os.Getenv("GO_ENV"), os.Getenv("ENVIRONMENT"), "development"),
		AttachStacktrace: true,
	})
	if err != nil {
		slog.Error("failed to initialize sentry", "error", err)
		return false
	}

	return true
}

func Flush() {
	sentry.Flush(flushTimeout)
}

func CaptureException(err error) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
}

func HTTPMiddleware() func(http.Handler) http.Handler {
	handler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,
		WaitForDelivery: false,
	})

	return func(next http.Handler) http.Handler {
		return handler.Handle(next)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
