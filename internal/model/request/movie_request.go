package request

// CreateMovieRequest representa el payload para crear una nueva película.
// swagger:model CreateMovieRequest
type CreateMovieRequest struct {
	// Título de la película
	// example: "El Gran Viaje"
	// required: true
	Title string `json:"title" binding:"required"`

	// Descripción de la película
	// example: "Una aventura épica"
	// required: true
	Description string `json:"description" binding:"required"`

	// Año de lanzamiento de la película
	// example: 1980
	// required: true
	Year int `json:"year" binding:"required"`

	// URL de la imagen de portada
	// example: "https://example.com/image.jpg"
	// required: true
	ImageURL string `json:"image_url" binding:"required"`

	// URL del video para streaming
	// example: "https://example.com/stream.mp4"
	// required: true
	StreamURL string `json:"stream_url" binding:"required"`

	// Género de la película
	// example: "Aventura"
	// required: true
	Genre string `json:"genre" binding:"required"`

	// Duración en minutos
	// example: 120
	// required: true
	Duration int `json:"duration" binding:"required"`
}

// UpdateMovieRequest representa el payload para actualizar una película.
// Los campos son punteros para permitir actualizaciones parciales.
// swagger:model UpdateMovieRequest
type UpdateMovieRequest struct {
	// Nuevo título de la película
	// example: "El Gran Viaje - Edición Especial"
	// required: false
	Title *string `json:"title,omitempty"`

	// Nueva descripción de la película
	// example: "Una aventura épica con escenas adicionales"
	// required: false
	Description *string `json:"description,omitempty"`

	// Nuevo año de lanzamiento
	// example: 1981
	// required: false
	Year *int `json:"year,omitempty"`

	// Nueva URL de la imagen de portada
	// example: "https://example.com/new_image.jpg"
	// required: false
	ImageURL *string `json:"image_url,omitempty"`

	// Nueva URL del video para streaming
	// example: "https://example.com/new_stream.mp4"
	// required: false
	StreamURL *string `json:"stream_url,omitempty"`

	// Nuevo género de la película
	// example: "Aventura - Edición Especial"
	// required: false
	Genre *string `json:"genre,omitempty"`

	// Nueva duración en minutos
	// example: 150
	// required: false
	Duration *int `json:"duration,omitempty"`
}
