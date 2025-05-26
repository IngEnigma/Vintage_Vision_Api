package request

// CreateProfileRequest representa el cuerpo de la solicitud para crear un perfil
// swagger:model CreateProfileRequest
type CreateProfileRequest struct {
	// Nombre del perfil (entre 1 y 10 caracteres)
	// required: true
	// example: "MiPerfil"
	// minLength: 1
	// maxLength: 10
	Name string `json:"name" binding:"required"`

	// URL válida del avatar del perfil
	// required: true
	// example: "https://ejemplo.com/avatar.jpg"
	// format: uri
	AvatarUrl string `json:"avatar_url" binding:"required,url"`
}

// UpdateProfileRequest representa los campos que se pueden actualizar en un perfil
// swagger:model UpdateProfileRequest
type UpdateProfileRequest struct {
	// Nuevo nombre del perfil (entre 1 y 10 caracteres, opcional)
	// example: "NuevoNombre"
	// minLength: 1
	// maxLength: 10
	Name *string `json:"name" binding:"omitempty,min=1,max=10"`

	// Nueva URL válida del avatar (opcional)
	// example: "https://ejemplo.com/nuevo-avatar.jpg"
	// format: uri
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url"`
}
