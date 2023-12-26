package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger = zap.NewNop()

var (
	_ http.ResponseWriter = (*loggedResponseWriter)(nil)
)

type info struct {
	size   int
	status int
}

type loggedResponseWriter struct {
	w http.ResponseWriter
	i *info
}

// Header - реализация метода Header интерфейса http.ResponseWriter
func (lrw *loggedResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

// Write - реализация метода Write интерфейса http.ResponseWriter
func (lrw *loggedResponseWriter) Write(body []byte) (int, error) {
	size, err := lrw.w.Write(body)
	lrw.i.size = size
	return size, err
}

// WriteHeader - реализация метода WriteHeader интерфейса http.ResponseWriter
func (lrw *loggedResponseWriter) WriteHeader(statusCode int) {
	lrw.i.status = statusCode
	lrw.w.WriteHeader(statusCode)
}

// Logger - метод доступа к логгеру
func Logger() *zap.Logger {
	return logger
}

func LogError(funcName string, err error) {
	Logger().Error("error",
		zap.String("func", "GetHandle"),
		zap.Error(err),
	)
}

// Init - метод инициализации логгера с уровнем по умолчанию
func Init(level string) error {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("error on parsing zap atomic level: %w", err)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atomicLevel
	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("error on building zap config: %w", err)
	}

	logger = zl
	return nil
}

// LoggerMiddleware - метод посредника для логирования данных запроса
func LoggerMiddleware(next http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		start := time.Now()

		loggerRW := loggedResponseWriter{
			w: w,
			i: &info{status: http.StatusOK},
		}
		next.ServeHTTP(&loggerRW, r)
		duration := time.Since(start)

		Logger().Info(">",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Int64("duration ms", duration.Milliseconds()),
		)

		Logger().Info("<",
			zap.Int("size", loggerRW.i.size),
			zap.Int("status", loggerRW.i.status),
		)
	}
	return http.HandlerFunc(hf)
}
