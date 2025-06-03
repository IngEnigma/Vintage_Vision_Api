package response

// MovieResponse representa la estructura de respuesta completa para una película.
// swagger:model MovieResponse
type MovieResponse struct {
	// ID único de la película
	// required: true
	// example: 1
	ID uint `json:"id"`

	// Título de la película
	// required: true
	// example: El Padrino
	Title string `json:"title"`

	// Descripción de la película
	// required: true
	// example: Una épica historia de crimen organizado en la América de los años 40
	Description string `json:"description"`

	// Año de lanzamiento
	// required: true
	// example: 1972
	Year int `json:"year"`

	// URL de la imagen de portada
	// required: true
	// example: https://cloudinary.com/example/padrino.jpg
	ImageURL string `json:"image_url"`

	// URL del video para streaming
	// required: true
	// example: https://cloudinary.com/example/padrino.mp4
	StreamURL string `json:"stream_url"`

	// Género de la película
	// required: true
	// example: Drama
	Genre string `json:"genre"`

	// Duración en minutos
	// required: true
	// example: 175
	Duration int `json:"duration"`
}

// MoviePreviewResponse representa una vista previa simplificada de película para listados.
// swagger:model MoviePreviewResponse
type MoviePreviewResponse struct {
	// ID único de la película en formato string
	// required: true
	// example: "1"
	ID string `json:"id"`

	// URL de la imagen miniatura
	// required: true
	// example: https://cloudinary.com/example/padrino_thumb.jpg
	ImageURL string `json:"image_url"`
}

// MovieDetailResponse representa los detalles extendidos de una película.
// swagger:model MovieDetailResponse
type MovieDetailResponse struct {
	// ID único de la película en formato string
	// required: true
	// example: "1"
	ID string `json:"id"`

	// Descripción extendida de la película
	// required: true
	// example: La saga de la familia Corleone, liderada por Vito Corleone y posteriormente por su hijo Michael.
	Description string `json:"description"`

	// Género principal de la película
	// required: true
	// example: Crimen
	Genre string `json:"genre"`

	// Año de lanzamiento
	// required: true
	// example: 1972
	Year int `json:"year"`

	// URL de la imagen en alta calidad
	// required: true
	// example: https://cloudinary.com/example/padrino_detail.jpg
	ImageURL string `json:"image_url"`
}

// MoviePlayerResponse representa los datos necesarios para el reproductor de video.
// swagger:model MoviePlayerResponse
type MoviePlayerResponse struct {
	// ID único de la película en formato string
	// required: true
	// example: "1"
	ID string `json:"id"`

	// Título de la película
	// required: true
	// example: El Padrino
	Title string `json:"title"`

	// URL del stream de video
	// required: true
	// example: https://cloudinary.com/example/padrino_stream.m3u8
	StreamURL string `json:"stream_url"`
}

// MovieListResponse representa una lista de películas con paginación.
// swagger:model MovieListResponse
type MovieTitleResponse struct {
	// Películas obtenidas
	// required: true
	ID string `json:"id"`

	// Título de la película
	// required: true
	Title string `json:"title"`
}
