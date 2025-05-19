package response

// AuthResponse define la estructura de respuesta para autenticación
// swagger:model AuthResponse
type AuthResponse struct {
	// Token JWT para autenticación
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	Token string `json:"token"`
}

// SuccessResponse define una respuesta genérica exitosa
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Mensaje de éxito
	// example: Operación realizada con éxito
	Message string `json:"message"`
}
