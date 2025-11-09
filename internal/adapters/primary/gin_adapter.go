package primary

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"juansecalvinio/go-auth-google/internal/core/domain"
	"juansecalvinio/go-auth-google/internal/ports/primary"

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
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas de Goth (Autenticación) - No usan directamente el userService
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/:provider", a.beginAuthHandler)
		authGroup.GET("/:provider/callback", a.completeAuthHandler)
	}

	// Rutas Protegidas - Usan el userService (Lógica de Negocio)
	router.GET("/user", a.userHandler)
	router.GET("/logout/:provider", a.logoutHandler)

	// Avatar proxy (evita problemas CORS/ORB y hotlinking)
	router.GET("/avatar", a.avatarHandler)
}

// avatarHandler proxya una imagen remota al cliente. Se aplica una whitelist
// de hosts para evitar SSRF y se devuelven cabeceras de cache.
func (a *GinAdapter) avatarHandler(c *gin.Context) {
	imgURL := c.Query("url")
	if imgURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing url parameter"})
		return
	}

	// Validar URL
	parsed, err := url.Parse(imgURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	// Whitelist de hosts (ajusta según necesites)
	host := parsed.Hostname()
	allowedHosts := []string{
		"lh3.googleusercontent.com",
		"googleusercontent.com",
		"avatars.githubusercontent.com",
	}
	ok := false
	for _, h := range allowedHosts {
		if strings.HasSuffix(host, h) {
			ok = true
			break
		}
	}
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "host not allowed"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build request"})
		return
	}
	// Agregar header User-Agent para evitar ciertos rechazos por hotlinking
	req.Header.Set("User-Agent", "go-auth-google-proxy/1.0")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "upstream returned non-200"})
		return
	}

	// Pasar Content-Type y copiar el body
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		c.Header("Content-Type", ct)
	} else {
		c.Header("Content-Type", "image/*")
	}
	c.Header("Cache-Control", "public, max-age=86400")
	c.Status(http.StatusOK)

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
		return
	}

}

// Handlers de Autenticación Goth
func (a *GinAdapter) beginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No provider specified"})
		return
	}

	// Añadir el provider a la query
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Configurar la sesión antes de comenzar el flujo de auth
	session, _ := gothic.Store.Get(c.Request, "gothic_session")
	session.Options.SameSite = http.SameSiteLaxMode // Importante para cookies cross-site
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	}
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (a *GinAdapter) completeAuthHandler(c *gin.Context) {
	// Obtener y validar el provider
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No provider specified"})
		return
	}

	// Añadir el provider a la query
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Intentar completar la autenticación
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Printf("ERROR - completeAuthHandler: %s - user data: %v\n", err.Error(), user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error completing authentication"})
		return
	}

	session, _ := gothic.Store.Get(c.Request, "gothic_session")
	session.Values["user_id"] = user.UserID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	session.Values["user_picture"] = user.AvatarURL

	if err := session.Save(c.Request, c.Writer); err != nil {
		fmt.Printf("ERROR - saving session: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving session"})
		return
	}

	redirectURL := c.Query("state")
	if redirectURL == "" {
		redirectURL = "http://localhost:5173/profile"
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

	fmt.Printf("User: %+v\n", session)

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

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
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
