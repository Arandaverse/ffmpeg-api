package routes

import (
	"github.com/gofiber/fiber/v2"
)

// IndexRoutes handles the index page route
type IndexRoutes struct{}

// NewIndexRoutes creates a new IndexRoutes instance
func NewIndexRoutes() *IndexRoutes {
	return &IndexRoutes{}
}

// Register registers all index routes
func (r *IndexRoutes) Register(router fiber.Router) {
	router.Get("/api/v1", r.handleIndex)
}

// handleIndex renders the index page
// @Summary Show index page
// @Description Display the main page with redirect button
// @Tags Index
// @Accept html
// @Produce html
// @Success 200 {string} string "HTML content"
// @Router /api/v1 [get]
func (r *IndexRoutes) handleIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "FFMPEG Serverless API",
	})
}
