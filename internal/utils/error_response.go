package utils

// ErrorResponse representa un error estándar de la API.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
