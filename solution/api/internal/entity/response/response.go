package response

// Error represents an error response message in JSON format.
type Error struct {
	Error string `json:"error"`
}

// Message represents a success response message in JSON format.
type Message struct {
	Message string `json:"message"`
}
