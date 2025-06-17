package infrastructure

import (
	infra_app "go-lti/internal/app"
	"go-lti/internal/lti"
	// "go-lti/lib/supabase"
	// "log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	infra_app.NewHttpHandler(v1)
	lti.NewHttpHandler(v1.Group("/lti"), ltiService)

	// sb := supabase.NewSupabaseClient()
	// if err := sb.Ping(); err != nil {
	// 	log.Fatalf("Supabase connection failed: %v", err)
	// }
	// log.Println("✅ Supabase connected successfully")

	app.Listen(":8080")
}
