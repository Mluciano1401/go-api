package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("el usuario no fue encontrado")
	ErrEmailAlreadyExists = errors.New("el email ya está en uso")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
)
