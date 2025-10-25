package secondary

import (
	domain "juansecalvinio/go-auth-google/internal/core/domain"
)

// AuthPort es la interfaz del puerto secundario (Autenticación/Goth).
// Define los métodos que la capa de servicio usará para interactuar con Goth.
type AuthPort interface {
	GetAuthenticatedUser(sessionID string) (*domain.User, error)
	Logout(sessionID string) error
}
