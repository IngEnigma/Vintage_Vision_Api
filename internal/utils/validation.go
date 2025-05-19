package utils

import (
	"net/url"
	"regexp"
	"unicode/utf8"
	"vintage-vision-api/internal/model/request"
)

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
