package response

// ErrorResponse representa una respuesta de error estándar de la API
// swagger:model
type ErrorResponse struct {
	// Código de estado HTTP
	// example: 400
	StatusCode int `json:"statusCode"`

	// Mensaje descriptivo del error
	// example: "Solicitud inválida"
	Message string `json:"message"`

	// Detalle técnico del error (opcional)
	// example: "Invalid email format"
	Error string `json:"error,omitempty"`
}
