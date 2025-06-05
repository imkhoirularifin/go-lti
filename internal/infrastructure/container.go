package infrastructure

import (
	"go-lti/internal/domain/interfaces"
	"go-lti/internal/lti"
	"go-lti/lib/config"
	"log"
)

var (
	cfg config.AppConfig

	ltiService interfaces.LtiService
)

func init() {
	var err error
	cfg, err = config.Setup()
	if err != nil {
		log.Fatalf("Failed to setup config: %v", err)
	}

	ltiService = lti.NewService(cfg)
}
