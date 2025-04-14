package errorz

// APIError represents a structured error for API responses.
type APIError struct {
	Status int    // HTTP status code
	Err    error  // Original error for debugging
	Msg    string // Human-readable error message for API response
}

// Error implements the error interface, returning the underlying error message.
func (e APIError) Error() string {
	return e.Err.Error()
}
