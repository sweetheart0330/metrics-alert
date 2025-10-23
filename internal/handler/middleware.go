package handler

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	compressReqHeader  = "Accept-Encoding"
	compressRespHeader = "Content-Encoding"
	compressFormat     = "gzip"
	hashSumHeader      = "Hashsha256"
	contentHeader      = "Content-Type"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (h Handler) CheckHashSum(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get(contentHeader)) == 0 || len(h.secretKey) == 0 {
			fmt.Printf("here1 %s, here: %d", r.Header.Get(contentHeader), len(h.secretKey))
			next.ServeHTTP(w, r)
			return
		}

		receivedSignature := r.Header.Get(hashSumHeader) // или как ты его назвал
		if receivedSignature == "" {
			fmt.Println("here2")
			http.Error(w, "Missing HMAC signature", http.StatusForbidden)
			return
		}

		fmt.Println("hash: ", receivedSignature)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		h := hmac.New(sha256.New, []byte(h.secretKey))
		_, err = h.Write(bodyBytes)
		if err != nil {
			http.Error(w, "Failed to compute HMAC", http.StatusInternalServerError)
			return
		}
		expectedSignature := h.Sum(nil)

		receivedBytes, err := hex.DecodeString(receivedSignature)
		if err != nil {
			http.Error(w, "Invalid signature format", http.StatusBadRequest)
			return
		}

		if !hmac.Equal(expectedSignature, receivedBytes) {
			http.Error(w, "Invalid HMAC signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h Handler) CompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get(compressReqHeader), compressFormat) {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (h Handler) DecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get(compressReqHeader), compressFormat) {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get(compressRespHeader), compressFormat) {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = struct {
				io.Reader
				io.Closer
			}{gz, r.Body}
		}

		next.ServeHTTP(w, r)
	})
}

func (h Handler) MiddlewareLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			reqID := middleware.GetReqID(r.Context())
			defer func() {
				h.log.Infow(
					"REQUEST COMPLETED",
					"reqID", reqID,
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"duration", time.Since(t1).String(),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
