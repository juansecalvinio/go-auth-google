package primary

import "juansecalvinio/go-auth-google/internal/core/domain"

// UserPort es la interfaz del puerto primario (la API/Web).
// Define los métodos que el Adaptador Gin llamará en la capa de servicio.
type UserPort interface {
	GetLoggedInUser(sessionID string) (*domain.User, error)
	LogoutUser(sessionID string) error
}
