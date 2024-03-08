package dto

type (
	Response struct {
		Success bool              `json:"success"`
		Message map[string]string `json:"message"`
		Data    interface{}       `json:"data,omitempty"`
		Meta    interface{}       `json:"meta,omitempty"`
		Error   *Error            `json:"error,omitempty"`
	}

	Error struct {
		Message string `json:"message,omitempty"`
	}
)
