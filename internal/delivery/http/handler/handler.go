package handler

import (
	"net/http"
	"strings"

	uc "github.com/ddalux/go_ai_fiber/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

// Handler wires HTTP routes to usecase
type Handler struct {
	uc      uc.UserUsecase
	swagger []byte
}

func NewHandler(u uc.UserUsecase, swagger []byte) *Handler {
	return &Handler{uc: u, swagger: swagger}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello World"})
	})
	app.Post("/register", h.register)
	app.Post("/login", h.login)
	app.Get("/me", h.me)
	app.Get("/swagger/doc.json", h.swaggerJSON)
	app.Get("/swagger", h.swaggerUI)
}

func (h *Handler) register(c *fiber.Ctx) error {
	var p struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Phone     string `json:"phone"`
		Birthday  string `json:"birthday"`
	}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := h.uc.Register(p.Email, p.Password, p.FirstName, p.LastName, p.Phone, p.Birthday); err != nil {
		if err == uc.ErrUserExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})
}

func (h *Handler) login(c *fiber.Ctx) error {
	var p struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	token, err := h.uc.Login(p.Email, p.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	return c.JSON(fiber.Map{"token": token})
}

func (h *Handler) me(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}
	token := strings.TrimPrefix(auth, prefix)
	user, err := h.uc.Me(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	return c.JSON(fiber.Map{
		"email":      user.Email,
		"firstname":  user.FirstName,
		"lastname":   user.LastName,
		"phone":      user.Phone,
		"birthday":   user.Birthday,
		"created_at": user.CreatedAt,
	})
}

func (h *Handler) swaggerJSON(c *fiber.Ctx) error {
	if len(h.swagger) == 0 {
		return c.Status(http.StatusInternalServerError).SendString("swagger not found")
	}
	c.Set("Content-Type", "application/json")
	return c.Send(h.swagger)
}

func (h *Handler) swaggerUI(c *fiber.Ctx) error {
	html := `<!doctype html><html><head><title>API Docs</title></head><body><redoc spec-url="/swagger/doc.json"></redoc><script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script></body></html>`
	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
