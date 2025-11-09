package main

import (
	"fmt"
	"juansecalvinio/go-auth-google/internal/adapters/primary"
	"juansecalvinio/go-auth-google/internal/adapters/secondary"
	"juansecalvinio/go-auth-google/internal/config"
	"juansecalvinio/go-auth-google/internal/core/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func initGoth(cfg *config.Config) *sessions.CookieStore {
	key := []byte(cfg.Security.SessionSecret)
	store := sessions.NewCookieStore(key)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 días
		HttpOnly: true,
		Secure:   cfg.Security.IsProd,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	goth.UseProviders(
		google.New(
			cfg.Auth.GoogleClientID,
			cfg.Auth.GoogleClientSecret,
			cfg.Auth.CallbackURL,
			cfg.Auth.Scopes...,
		),
	)

	return store
}

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Inicializar Goth y obtener store
	store := initGoth(cfg)

	// --- 1. CONFIGURACIÓN E INFRAESTRUCTURA (Adaptadores y Puertos) ---

	// Adaptador de Autenticación (Puerto Secundario - Goth)
	gothAdapter := secondary.NewGothAdapter(store)

	// --- 2. CORE (Dominio y Servicio) ---

	// Servicio de Usuario (Lógica de Negocio)
	userService := service.NewUserService(gothAdapter)

	// --- 3. ADAPTADOR WEB (Puerto Primario - Gin) ---

	// Adaptador Gin para la API
	ginAdapter := primary.NewGinAdapter(userService)

	// --- 4. INICIO DEL SERVIDOR ---
	router := gin.Default()
	router.Use(primary.CORSMiddleware())

	ginAdapter.RegisterRoutes(router)

	fmt.Println("Backend Hexagonal (Gin/Goth) corriendo en :8080")
	log.Fatal(router.Run(":8080"))
}
