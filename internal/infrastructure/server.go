package infrastructure

import (
	infra_app "go-lti/internal/app"
	"go-lti/internal/canvas"
	"go-lti/internal/lti"
	"go-lti/lib/common"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func Run() {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: common.ErrorHandler,
		},
	)

	app.Use(recover.New())
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	infra_app.NewHttpHandler(v1)
	lti.NewHttpHandler(v1.Group("/lti"), ltiService)
	canvas.NewHttpHandler(v1.Group("/canvas"), canvasService)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Error().Err(err).Msg("Failed to start server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info().Msg("Shutting down server")
	app.Shutdown()
	log.Info().Msg("Running cleanup tasks")

	// Your cleanup tasks here
	// db.Close()
	// redisConn.Close()

	log.Info().Msg("Server shutdown complete")
}
