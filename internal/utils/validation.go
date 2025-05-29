package utils

import (
	"net/url"
	"regexp"
	"unicode/utf8"
	"vintage-vision-api/internal/model/request"

	"vintage-vision-api/internal/constants"
)

var validGenres = map[string]bool{
	constants.CategoryComedy:          true,
	constants.CategoryDrama:           true,
	constants.CategoryHorror:          true,
	constants.CategorySciFi:           true,
	constants.CategoryMusical:         true,
	constants.CategoryPropaganda:      true,
	constants.CategoryDibujosAnimados: true,
}

func IsValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return utf8.RuneCountInString(password) >= 8
}

func IsNonEmpty(s string) bool {
	return len(s) > 0 && len(regexp.MustCompile(`\S`).FindString(s)) > 0
}

func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsValidProfileName(name string) bool {
	if len(name) == 0 || len(name) > 10 {
		return false
	}
	return true
}

func ValidatePaginationParams(page, limit int) (int, int) {
	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 {
		limit = constants.DefaultLimit
	} else if limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}
	return page, limit
}

func ValidateRegisterInput(req request.RegisterRequest) map[string]string {
	errors := map[string]string{}

	if !IsNonEmpty(req.Email) || !IsValidEmail(req.Email) {
		errors["email"] = "Email inválido"
	}
	if !IsValidPassword(req.Password) {
		errors["password"] = "La contraseña debe tener al menos 8 caracteres"
	}

	return errors
}

func ValidateLoginInput(req request.LoginRequest) map[string]string {
	errors := map[string]string{}

	if !IsNonEmpty(req.Email) || !IsValidEmail(req.Email) {
		errors["email"] = "Email inválido"
	}
	if !IsNonEmpty(req.Password) {
		errors["password"] = "La contraseña no puede estar vacía"
	}

	return errors
}

func ValidateCreateMovieInput(req request.CreateMovieRequest) map[string]string {
	errors := map[string]string{}

	if !IsNonEmpty(req.Title) {
		errors["title"] = "El título es obligatorio"
	}
	if !IsNonEmpty(req.Description) {
		errors["description"] = "La descripción es obligatoria"
	}
	if req.Year < 1888 || req.Year > 2000 {
		errors["year"] = "El año debe ser válido"
	}
	if !IsValidURL(req.ImageURL) {
		errors["image_url"] = "La URL de la imagen no es válida"
	}
	if !IsValidURL(req.StreamURL) {
		errors["stream_url"] = "La URL de streaming no es válida"
	}
	if !IsNonEmpty(req.Genre) {
		errors["genre"] = "El género es obligatorio"
	}
	if req.Duration <= 0 {
		errors["duration"] = "La duración debe ser mayor a 0 minutos"
	}

	return errors
}

func ValidateUpdateMovieInput(req request.UpdateMovieRequest) map[string]string {
	errors := map[string]string{}

	if req.Title != nil && !IsNonEmpty(*req.Title) {
		errors["title"] = "El título no puede estar vacío"
	}
	if req.Description != nil && !IsNonEmpty(*req.Description) {
		errors["description"] = "La descripción no puede estar vacía"
	}
	if req.Year != nil && (*req.Year < 1888 || *req.Year > 2100) {
		errors["year"] = "El año debe ser válido"
	}
	if req.ImageURL != nil && !IsValidURL(*req.ImageURL) {
		errors["image_url"] = "La URL de la imagen no es válida"
	}
	if req.StreamURL != nil && !IsValidURL(*req.StreamURL) {
		errors["stream_url"] = "La URL de streaming no es válida"
	}
	if req.Genre != nil && !IsNonEmpty(*req.Genre) {
		errors["genre"] = "El género no puede estar vacío"
	}
	if req.Duration != nil && *req.Duration <= 0 {
		errors["duration"] = "La duración debe ser mayor a 0 minutos"
	}

	return errors
}

func BuildMovieUpdateMap(req request.UpdateMovieRequest) map[string]interface{} {
	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.StreamURL != nil {
		updates["stream_url"] = *req.StreamURL
	}
	if req.Genre != nil {
		updates["genre"] = *req.Genre
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}

	return updates
}

func IsValidGenre(category string) bool {
	return validGenres[category]
}
