package middleware

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func JSONResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &ResponseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		if wrappedWriter.statusCode >= 200 && wrappedWriter.statusCode < 300 {
			return
		}

		// Если нет данных, отправляем статус с пустым объектом
		resp := Response{
			StatusCode: wrappedWriter.statusCode,
			Data:       map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		StatusCode: statusCode,
		Data:       data,
	}

	json.NewEncoder(w).Encode(resp)
}
