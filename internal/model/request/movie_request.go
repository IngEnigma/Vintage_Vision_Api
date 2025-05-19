package request

// CreateMovieRequest representa el payload para crear una nueva película.
// swagger:model CreateMovieRequest
type CreateMovieRequest struct {
	// Título de la película
	// required: true
	Title string `json:"title"`

	// Descripción de la película
	// required: true
	Description string `json:"description"`

	// Año de lanzamiento de la película
	// required: true
	Year int `json:"year"`

	// URL de la imagen de portada
	// required: true
	ImageURL string `json:"image_url"`

	// URL del video para streaming
	// required: true
	StreamURL string `json:"stream_url"`

	// Género de la película
	// required: true
	Genre string `json:"genre"`

	// Duración en minutos
	// required: true
	Duration int `json:"duration"`
}

// UpdateMovieRequest representa el payload para actualizar una película.
// Los campos son punteros para permitir actualizaciones parciales.
// swagger:model UpdateMovieRequest
type UpdateMovieRequest struct {
	// Nuevo título de la película
	// required: false
	Title *string `json:"title,omitempty"`

	// Nueva descripción de la película
	// required: false
	Description *string `json:"description,omitempty"`

	// Nuevo año de lanzamiento
	// required: false
	Year *int `json:"year,omitempty"`

	// Nueva URL de la imagen de portada
	// required: false
	ImageURL *string `json:"image_url,omitempty"`

	// Nueva URL del video para streaming
	// required: false
	StreamURL *string `json:"stream_url,omitempty"`

	// Nuevo género de la película
	// required: false
	Genre *string `json:"genre,omitempty"`

	// Nueva duración en minutos
	// required: false
	Duration *int `json:"duration,omitempty"`
}
