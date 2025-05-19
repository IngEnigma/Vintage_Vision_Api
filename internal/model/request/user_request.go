package request

// RegisterRequest define la estructura para el registro de usuarios
// swagger:model RegisterRequest
type RegisterRequest struct {
	// Email del usuario (debe ser válido)
	// required: true
	// example: usuario@ejemplo.com
	Email string `json:"email" binding:"required,email"`

	// Contraseña (mínimo 8 caracteres)
	// required: true
	// example: Password123!
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest define la estructura para el inicio de sesión
// swagger:model LoginRequest
type LoginRequest struct {
	// Email del usuario
	// required: true
	// example: usuario@ejemplo.com
	Email string `json:"email" binding:"required,email"`

	// Contraseña del usuario
	// required: true
	// example: Password123!
	Password string `json:"password" binding:"required"`
}
