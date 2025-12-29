package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ahrdadan/image-metadata-viewer/src/internal/handlers"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Initialize template engine
	engine := html.New("./src/web/templates", ".html")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views:                 engine,
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          15 * time.Second,
		BodyLimit:             21 * 1024 * 1024, // 21MB
		DisableStartupMessage: false,
		AppName:               "Image Metadata Viewer v2.0",
		ErrorHandler:          customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Serve static files (CSS, JS, images)
	app.Static("/static", "./src/web/static")

	// Initialize services
	imageService := services.NewImageService()
	blobStore := services.NewBlobStore(time.Hour)

	// Initialize handlers
	webHandler := handlers.NewWebHandler(imageService, blobStore)
	apiHandler := handlers.NewAPIHandler(imageService)

	// API routes
	api := app.Group("/api")
	api.Get("/*", apiHandler.HandleGetMetadata)
	api.Post("/", apiHandler.HandlePostMetadata)

	// Web routes
	app.Get("/", webHandler.HandleHome)
	app.Get("/docs", webHandler.HandleDocs)
	app.Get("/go", webHandler.HandleForm)
	app.Post("/upload", webHandler.HandleUpload)
	app.Get("/blob/:id", webHandler.HandleBlob)
	app.Get("/*", webHandler.HandleView)

	// Get port from environment or use default
	port := getPort()

	log.Printf("🚀 Server starting on http://localhost%s", port)
	log.Printf("📚 API documentation available at /api")
	log.Printf("🎨 Web interface available at /")

	if err := app.Listen(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getPort returns the port from environment or default
func getPort() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return ":8080"
	}
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}
	return port
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Check if request is API
	if strings.HasPrefix(c.Path(), "/api") {
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"error":   message,
		})
	}

	// Web error page
	return c.Status(code).SendString(message)
}
