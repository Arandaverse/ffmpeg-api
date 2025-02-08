package response

import (
	"encoding/json"
	"ffmpeg-api/internal/logger"
	"ffmpeg-api/internal/validation"
	"net/http"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

// ResponseWriter wraps http.ResponseWriter with utility methods
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

// WriteHeader captures the status code and passes it to the underlying ResponseWriter
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the HTTP status code of the response
func (w *ResponseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// writeJSON writes a JSON response with headers
func (w *ResponseWriter) writeJSON(statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// isSuccess checks if the status code represents success
func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	rw := NewResponseWriter(w)
	resp := Response{
		Success: isSuccess(statusCode),
		Data:    data,
	}
	rw.writeJSON(statusCode, resp)
}

// Error sends a JSON error response with the given status code and error details
func Error(w http.ResponseWriter, statusCode int, code string, message interface{}) {
	rw := NewResponseWriter(w)
	resp := Response{
		Success: false,
		Error: &APIError{
			Type:    code,
			Message: message,
		},
	}
	rw.writeJSON(statusCode, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), message)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), message)
}

// Conflict sends a 409 Conflict response
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, http.StatusText(http.StatusConflict), message)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), message)
}

// ValidationError sends a 400 Bad Request response with validation details
func ValidationError(w http.ResponseWriter, err error) {
	if validationErrors, ok := err.(*validation.ValidationErrors); ok {
		Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), validationErrors.Errors)
		return
	}
	Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
}

// InvalidCredentials sends a 401 Unauthorized response for invalid credentials
func InvalidCredentials(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid credentials")
}
