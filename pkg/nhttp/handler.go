package nhttp

import (
	"encoding/json"
	"github.com/diarikom/running-app/running-app-api/pkg/nlog"
	"html"
	"net/http"
	"time"
)

const (
	KeyContentType   = "Content-Type"
	KeyAuthorization = "Authorization"
	ContentTypeJSON  = "application/json; charset=utf-8"
	ContentTypeXML   = "application/xml; charset=utf-8"
	ContentTypeHTML  = "text/html; charset=utf-8"
)

type HandlerFunc func(*http.Request) (*Success, error)

type Handler struct {
	Fn     HandlerFunc
	Logger nlog.Logger
}

// ServeHTTP implement http.Handler interface to write success or error response in JSON
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Init start
	start := time.Now()

	// Init http status
	var httpStatus int

	// Execute handler
	result, err := h.Fn(r)

	// If an error returned, return error
	if err != nil {
		httpStatus = h.sendErrorJSON(w, err)
	} else if result == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		// if header exist, add header to response
		if result.Header != nil {
			for k, v := range result.Header {
				w.Header().Set(k, v)
			}
		}
		// send json success
		httpStatus = h.sendJSON(w, http.StatusOK, result)
	}

	// Log elapsed time
	h.Logger.Infof("HTTP Status: %d, Request: %s %s, Time elapsed: %s", httpStatus, r.Method,
		html.EscapeString(r.URL.Path), time.Since(start))
}

// sendJSON write response in JSON
func (h Handler) sendJSON(w http.ResponseWriter, httpStatus int, obj interface{}) int {
	// Add content type
	w.Header().Add(KeyContentType, ContentTypeJSON)
	// Write http status
	w.WriteHeader(httpStatus)
	// Send JSON response
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {

	}
	// Return httpStatus
	return httpStatus
}

// sendErrorJSON write error response in JSON
func (h Handler) sendErrorJSON(w http.ResponseWriter, err error) int {
	// CastError error to Error
	apiError := CastError(err)
	// Send error json
	return h.sendJSON(w, apiError.Status, apiError)
}

// parseJSON parse json request body to o (target) and returns error
func ParseJSON(o interface{}, r *http.Request) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(o); err != nil {
		return err
	}
	return nil
}
