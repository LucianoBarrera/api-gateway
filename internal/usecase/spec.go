package usecase

import (
	"net/http"
)

// RequestForwarder interface for mocking in tests
type RequestForwarder interface {
	ForwardRequest(w http.ResponseWriter, r *http.Request, serviceName string)
}
