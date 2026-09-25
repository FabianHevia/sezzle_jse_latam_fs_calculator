package httpapi

// RequestDTO defines the incoming payload for arithmetic endpoints.
// Pointer fields allow distinguishing between omitted fields and explicit 0 values.
type RequestDTO struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

// SuccessDTO represents a successful 200 response payload.
type SuccessDTO struct {
	Result float64 `json:"result"`
}

// ErrorEnvelope represents the standard error response wrapper.
type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains machine-readable error codes and human-readable messages.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HealthDTO represents the healthcheck response structure.
type HealthDTO struct {
	Status string `json:"status"`
}
