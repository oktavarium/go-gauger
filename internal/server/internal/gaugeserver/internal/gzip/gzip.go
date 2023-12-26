package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressedWriter struct {
	w          http.ResponseWriter
	gzipWriter *gzip.Writer
}

var (
	_ http.ResponseWriter = (*compressedWriter)(nil)
)

func newGzipWriter(w http.ResponseWriter) *compressedWriter {
	gzipWriter := gzip.NewWriter(w)
	w.Header().Set("Content-Encoding", "gzip")
	return &compressedWriter{w, gzipWriter}
}

func (c *compressedWriter) Close() {
	c.gzipWriter.Close()
}

func (c *compressedWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressedWriter) Write(data []byte) (int, error) {
	return c.gzipWriter.Write(data)
}

func (c *compressedWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

type compressedReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressedReader(r io.ReadCloser) (*compressedReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("error on creating gzipReader: %w", err)
	}

	return &compressedReader{r, zr}, nil
}

func (c *compressedReader) Read(data []byte) (int, error) {
	return c.zr.Read(data)
}

func (c *compressedReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf(
			"error on closing reader in compressed reader: %w",
			err,
		)
	}
	return c.zr.Close()
}

// GzipMiddleware - метод посредника для архивирования и разархивирования запросов
func GzipMiddleware(next http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		// сохраняем оригинальный врайтер
		ow := w

		// проверяем то, что можем отправлять клиенту
		contentEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			gzipWriter := newGzipWriter(w)
			ow = gzipWriter
			defer gzipWriter.Close()
		}

		// проверяем, что можем принимать от клиента
		acceptEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			gzipReader, err := newCompressedReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gzipReader
			defer gzipReader.Close()
		}
		next.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(hf)
}
