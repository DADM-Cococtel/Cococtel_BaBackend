package defines

import "errors"

var (
	// Errores generales
	ErrShouldBindJSON = errors.New("Error obteniendo los datos")

	// Errores en el registro de usuario
	ErrAlreadyExists = errors.New("El usuario ya se encuentra registrado en la base de datos")

	// Error en el login
	ErrIncorrectPassword = errors.New("Contraseña incorrecta")
	ErrInvalidOTP        = errors.New("Codigo de seguridad invalido")

	ErrAuthAlreadyEnabled = errors.New("Ya se encuentra habilitado el login de dos factores")

	// Authorization
	ErrNotFoundAuthHeader = errors.New("Missing authorization header: x-auth-token")
	ErrInvalidToken       = errors.New("Token invalido")
)
