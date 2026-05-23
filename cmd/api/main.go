package main

import (
	"fmt"
	"log"

	"github.com/Mluciano1401/go-api/config"
	"github.com/Mluciano1401/go-api/internal/application"
	"github.com/Mluciano1401/go-api/internal/infrastructure/persistence"
	"github.com/Mluciano1401/go-api/internal/infrastructure/security"
	httpiface "github.com/Mluciano1401/go-api/internal/interfaces/http"
	"github.com/Mluciano1401/go-api/internal/interfaces/http/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	log.Printf("DB: host=%s port=%s user=%s pass=%s name=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}

	if err := db.AutoMigrate(&persistence.UserModelForMigration{}); err != nil {
		log.Fatalf("error en migración: %v", err)
	}

	userRepo := persistence.NewUserRepository(db)
	hasher := security.NewBcryptHasher()
	jwtGen := security.NewJWTGenerator(cfg.JWTSecret, cfg.JWTExpiresHrs)

	authService := application.NewAuthService(userRepo, hasher, jwtGen)
	userService := application.NewUserService(userRepo, hasher)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	router := httpiface.NewRouter(authHandler, userHandler, cfg.JWTSecret)
	addr := ":" + cfg.ServerPort
	log.Printf("Servidor escuchando en %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("error arrancando el servidor: %v", err)
	}
}
