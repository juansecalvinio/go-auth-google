package primary

import (
	"context"
	"fmt"
	"juansecalvinio/go-auth-google/internal/core/domain"
	"juansecalvinio/go-auth-google/internal/ports/primary"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

// contextKey es un tipo personalizado para las claves del contexto
type contextKey string

// Definición de claves del contexto
const (
	ProviderKey contextKey = "provider"
	StateKey    contextKey = "state"
)

// GinAdapter es el Adaptador que expone la API web usando Gin.
type GinAdapter struct {
	userService primary.UserPort
}

// NewGinAdapter crea una nueva instancia y establece la dependencia.
func NewGinAdapter(userPort primary.UserPort) *GinAdapter {
	return &GinAdapter{userService: userPort}
}

// RegisterRoutes configura todas las rutas de la API.
func (a *GinAdapter) RegisterRoutes(router *gin.Engine) {
	// Rutas de Goth (Autenticación) - No usan directamente el userService
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/:provider", a.beginAuthHandler)
		authGroup.GET("/:provider/callback", a.completeAuthHandler)
	}

	// Rutas Protegidas - Usan el userService (Lógica de Negocio)
	router.GET("/user", a.userHandler)
	router.GET("/logout/:provider", a.logoutHandler)
}

// Handlers de Autenticación Goth
func (a *GinAdapter) beginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	ctx := context.WithValue(c.Request.Context(), ProviderKey, provider)

	c.Request = c.Request.WithContext(ctx)

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *GinAdapter) completeAuthHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Printf("ERROR - completeAuthHandler: %s - user data: %s", err.Error(), user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error complete user auth"})
		return
	}

	session, _ := gothic.Store.Get(c.Request, "gothic_session")
	session.Values["user_id"] = user.UserID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	session.Values["user_picture"] = user.AvatarURL

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Printf("ERROR - completeAuthHandler: %s - user data: %s", err.Error(), user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error complete user auth"})
		return
	}

	redirectURL, ok := c.Request.Context().Value(StateKey).(string)
	if !ok {
		redirectURL = "http://localhost:5173" // URL por defecto
	}
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// Handlers Protegidos (Lógica de Negocio)

// userHandler ahora utiliza el GIN ADAPTER para extraer la sesión
// y crear un modelo de DOMINIO para enviarlo al frontend.
func (a *GinAdapter) userHandler(c *gin.Context) {
	session, err := gothic.Store.Get(c.Request, "gothic_session")
	if err != nil || session.Values["user_id"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado o sesión expirada."})
		return
	}

	// Extraer datos de la sesión (detalle de infraestructura) y mapear a la entidad de Dominio.
	// Esto es un **Anti-Corrupción Layer (ACL)** simplificado.
	user := &domain.User{
		ID:      fmt.Sprintf("%s", session.Values["user_id"]),
		Name:    fmt.Sprintf("%s", session.Values["user_name"]),
		Email:   fmt.Sprintf("%s", session.Values["user_email"]),
		Picture: fmt.Sprintf("%s", session.Values["user_picture"]),
	}

	// En un caso más complejo, llamaríamos a:
	// loggedInUser, _ := a.userService.GetLoggedInUser(sessionID)
	c.JSON(http.StatusOK, user)
}

func (a *GinAdapter) logoutHandler(c *gin.Context) {
	// El Gin Adapter maneja la limpieza de la sesión (detalle de infraestructura).
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error handling logout"})
	}
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/")
}

// Configuración de CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
