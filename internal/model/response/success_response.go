package response

// SuccessResponse define una respuesta genérica exitosa
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Indica si la operación fue exitosa
	// example: true
	Success bool `json:"success"`

	// Mensaje descriptivo del resultado
	// example: Operación realizada con éxito
	Message string `json:"message"`

	// Datos adicionales de la operación (opcional)
	Data interface{} `json:"data,omitempty"`
}
