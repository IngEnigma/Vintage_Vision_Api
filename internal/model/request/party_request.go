package request

// CreatePartyRequest representa el payload para crear una nueva party.
// swagger:model CreatePartyRequest
type CreatePartyRequest struct {
	// ID del usuario que será el anfitrión de la party
	// example: 1
	// required: true
	HostID uint `json:"host_id" binding:"required"`

	// ID de la película que se verá en la party
	// example: 42
	// required: true
	MovieID uint `json:"movie_id" binding:"required"`
}

// JoinPartyRequest representa el payload para unirse a una party existente.
// swagger:model JoinPartyRequest
type JoinPartyRequest struct {
	// Código de la party a la que se quiere unir
	// example: "ABCD12"
	// required: true
	Code string `json:"code" binding:"required"`

	// ID del perfil que se unirá a la party
	// example: 3
	// required: true
	ProfileID uint `json:"profile_id" binding:"required"`
}
