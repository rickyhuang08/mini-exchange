package entity

type APIResponse struct {
	Status  string      `json:"status"`            // "success" or "error"
	Message string      `json:"message,omitempty"` // Optional message
	Data    interface{} `json:"data,omitempty"`    // Any payload (e.g. LoginResponse)
	Error   string      `json:"error,omitempty"`   // Error details if status == "error"
}