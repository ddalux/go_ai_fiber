package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddalux/go_ai_fiber/internal/delivery/http/handler"
	repo "github.com/ddalux/go_ai_fiber/internal/repository"
	uc "github.com/ddalux/go_ai_fiber/internal/usecase"
)

// Swagger will be loaded from file at runtime

func main() {
	// setup DB
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.AutoMigrate(&repo.User{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	r := repo.NewGormUserRepo(db)
	u := uc.NewUserUsecase(r, []byte(getJWTSecret()))

	// load swagger.json
	swaggerBytes, err := os.ReadFile("swagger.json")
	if err != nil {
		log.Printf("warning: failed to read swagger.json: %v", err)
		swaggerBytes = []byte("{}")
	}

	app := fiber.New()
	app.Use(logger.New())

	h := handler.NewHandler(u, swaggerBytes)
	h.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "please-change-this-secret"
}
