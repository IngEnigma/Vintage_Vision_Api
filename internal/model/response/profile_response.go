package response

// ProfileResponse representa los datos de un perfil devuelto en la respuesta
type ProfileResponse struct {
	// ID del perfil
	ID uint `json:"id" example:"1"`
	// Nombre del perfil
	Name string `json:"name" example:"John"`
	// URL del avatar
	AvatarUrl string `json:"avatar_url" example:"https://example.com/avatar.png"`
}
