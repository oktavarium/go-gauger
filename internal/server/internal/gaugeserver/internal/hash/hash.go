package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

const HashHeader = "Hashsha256"

type hashedWriter struct {
	http.ResponseWriter

	mw io.Writer
	b  *bytes.Buffer
	k  []byte
}

func newHashedWriter(w http.ResponseWriter, key []byte) *hashedWriter {
	b := new(bytes.Buffer)
	mw := io.MultiWriter(w, b)
	return &hashedWriter{w, mw, b, key}
}

func (h *hashedWriter) Write(data []byte) (int, error) {
	return h.mw.Write(data)
}

func (h *hashedWriter) hash() string {
	return hashData(h.k, h.b.Bytes())
}

func hashData(key []byte, data []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hashedData := mac.Sum(nil)
	return hex.EncodeToString(hashedData)
}

func checkHash(key []byte, data []byte, hash string) error {
	clientHash := hashData(key, data)
	if clientHash != hash {
		return fmt.Errorf("hashes are not equal")
	}
	return nil
}

// HashMiddleware - посредник для проверки подлинности клиента
func HashMiddleware(key []byte) func(http.Handler) http.Handler {
	nextF := func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := r.Header[HashHeader]; ok {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Logger().Error("error",
						zap.String("func", "HashMiddleware"),
						zap.Error(err))

					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(bytes.NewReader(body))

				err = checkHash(key, body, r.Header.Get(HashHeader))
				if err != nil {
					logger.Logger().Error("error",
						zap.String("func", "HashMiddleware"),
						zap.Error(err))

					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			hashedWriter := newHashedWriter(w, key)

			next.ServeHTTP(hashedWriter, r)

			hash := hashedWriter.hash()
			hashedWriter.Header().Set(HashHeader, hash)
		}
		return http.HandlerFunc(hf)
	}

	return nextF
}
