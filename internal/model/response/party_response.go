package response

// CreatePartyResponse representa la respuesta al crear una nueva fiesta
// swagger:model CreatePartyResponse
type CreatePartyResponse struct {
	// Success indica si la operación fue exitosa
	Message string `json:"message"`
	// PartyCode es el código único de la fiesta creada
	// example: "1k36v2"
	// Length: 6
	PartyCode string `json:"party_code"`
}
