package service

import (
	"juansecalvinio/go-auth-google/internal/core/domain"
	"juansecalvinio/go-auth-google/internal/ports/secondary"
)

// UserService implementa la lógica de negocio.
type UserService struct {
	authPort secondary.AuthPort
}

// NewUserService crea una nueva instancia del servicio.
func NewUserService(authPort secondary.AuthPort) *UserService {
	return &UserService{authPort: authPort}
}

// GetLoggedInUser ejecuta la lógica: obtener los datos del usuario logueado.
func (s *UserService) GetLoggedInUser(sessionID string) (*domain.User, error) {
	// La lógica de negocio delega la recuperación al Adaptador de Autenticación.
	return s.authPort.GetAuthenticatedUser(sessionID)
}

// LogoutUser ejecuta la lógica: limpiar la sesión.
func (s *UserService) LogoutUser(sessionID string) error {
	return s.authPort.Logout(sessionID)
}
