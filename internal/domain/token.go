package domain

type TokenGenerator interface {
	Generate(userID uint, isAdmin bool) (string, error)
}
