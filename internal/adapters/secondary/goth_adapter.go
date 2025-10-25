package secondary

import (
	"fmt"
	"juansecalvinio/go-auth-google/internal/core/domain"
	"juansecalvinio/go-auth-google/internal/ports/secondary"

	"github.com/gorilla/sessions"
)

// GothAdapter es el Adaptador para la autenticación externa (Goth/Google).
type GothAdapter struct {
	store sessions.Store
}

// NewGothAdapter crea una nueva instancia del adaptador.
func NewGothAdapter(store sessions.Store) secondary.AuthPort {
	return &GothAdapter{store: store}
}

// GetAuthenticatedUser implementa el método del puerto.
// Nota: En una arquitectura real, pasarías el *http.Request al Adaptador,
// pero aquí simplificamos extrayendo directamente de la sesión (usando un SessionID ficticio).
func (a *GothAdapter) GetAuthenticatedUser(sessionID string) (*domain.User, error) {
	return nil, fmt.Errorf("La lectura de la sesión se maneja directamente en el adaptador Gin por la dependencia de gothic")
}

func (a *GothAdapter) Logout(sessionID string) error {
	return fmt.Errorf("La función Logout se maneja directamente en el adaptador Gin por la dependencia de Gothic")
}
