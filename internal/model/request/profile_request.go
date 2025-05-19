package request

// CreateProfileRequest representa el cuerpo de la solicitud para crear un perfil
type CreateProfileRequest struct {
	// Nombre del perfil
	Name string `json:"name" binding:"required" example:"John"`
	// URL del avatar (debe ser una URL válida)
	AvatarUrl string `json:"avatar_url" binding:"required,url" example:"https://example.com/avatar.png"`
}

// UpdateProfileRequest representa los campos que se pueden actualizar en un perfil
type UpdateProfileRequest struct {
	// Nuevo nombre del perfil (opcional)
	Name *string `json:"name" binding:"omitempty,min=1,max=10" example:"Jane"`
	// Nueva URL del avatar (opcional)
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url" example:"https://example.com/avatar2.png"`
}
