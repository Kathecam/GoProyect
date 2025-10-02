package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/Kathecam/go-tasks-api/internal/config"
	"github.com/Kathecam/go-tasks-api/internal/handlers"
	"github.com/Kathecam/go-tasks-api/internal/middleware"
)

func main() {
	// Cargar .env en desarrollo
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Crear instancia de Fiber con configuración
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName + " v" + cfg.AppVersion,
		ServerHeader: "Fiber",
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	// Logging condicional basado en DEBUG
	if cfg.Debug {
		log.Println("Running in debug mode")
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	} else {
		app.Use(logger.New())
	}

	// Middlewares globales (orden importa)
	app.Use(recover.New())             // Captura panics
	app.Use(cors.New())                // CORS básico
	app.Use(middleware.ErrorHandler()) // Manejo de errores centralizado

	// Crear handlers
	taskHandler := handlers.NewTaskHandler()

	// Rutas de sistema (fuera de versionado)
	app.Get("/health/ping", healthHandler(cfg))

	// Debug endpoint (solo en desarrollo)
	if cfg.Debug {
		app.Get("/debug/config", debugConfigHandler(cfg))
	}

	// Grupo API v1
	v1 := app.Group("/api/v1")
	v1.Use(versionMiddleware(cfg)) // Middleware para agregar versión en headers

	v1.Get("/version", versionHandler(cfg))
	// Grupo para tareas con nuevos handlers
	tasks := v1.Group("/tasks")
	tasks.Get("/", taskHandler.GetTasks)         // GET /api/v1/tasks
	tasks.Post("/", taskHandler.CreateTask)      // POST /api/v1/tasks
	tasks.Get("/:id", taskHandler.GetTaskByID)   // GET /api/v1/tasks/:id
	tasks.Put("/:id", taskHandler.UpdateTask)    // PUT /api/v1/tasks/:id
	tasks.Delete("/:id", taskHandler.DeleteTask) // DELETE /api/v1/tasks/:id

	// Iniciar servidor
	address := cfg.Host + ":" + cfg.Port
	log.Printf("Starting server on %s in %s mode", address, cfg.Environment)
	log.Fatal(app.Listen(address))
}

// Handler simple que retorna JSON
func healthHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   cfg.AppName,
			"message":   "Pong! 🔥",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"version":   cfg.AppVersion,
		})
	}
}

func versionHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		buildTime := "development"
		if cfg.IsProduction() {
			buildTime = "production"
		}

		return c.JSON(fiber.Map{
			"version":     cfg.AppVersion,
			"build_time":  buildTime,
			"environment": cfg.Environment,
		})
	}
}

func versionMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-API-Version", cfg.AppVersion)
		return c.Next()
	}
}

func debugConfigHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mostrar config SIN secrets
		return c.JSON(fiber.Map{
			"app_name":      cfg.AppName,
			"app_version":   cfg.AppVersion,
			"environment":   cfg.Environment,
			"debug":         cfg.Debug,
			"host":          cfg.Host,
			"port":          cfg.Port,
			"read_timeout":  cfg.ReadTimeout.String(),
			"write_timeout": cfg.WriteTimeout.String(),
			// NO incluir: JWT_SECRET, DATABASE_URL
		})
	}
}
