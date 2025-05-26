package response

// ProfileResponse representa los datos de un perfil devuelto en la respuesta
// swagger:model ProfileResponse
type ProfileResponse struct {
	// ID único del perfil
	// example: 1
	ID uint `json:"id"`

	// Nombre del perfil (máximo 10 caracteres)
	// example: "MiPerfil"
	// maxLength: 10
	Name string `json:"name"`

	// URL completa del avatar del perfil
	// example: "https://ejemplo.com/avatar.jpg"
	// format: uri
	AvatarUrl string `json:"avatar_url"`
}
