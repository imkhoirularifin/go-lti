package infrastructure

import (
	"go-lti/internal/canvas"
	"go-lti/internal/domain/interfaces"
	"go-lti/internal/lti"
	"go-lti/lib/config"
	"go-lti/lib/httpclient"
	"log"
	"time"
)

var (
	cfg config.AppConfig

	httpClient httpclient.HttpClient

	ltiService    interfaces.LtiService
	canvasService interfaces.CanvasService
)

func init() {
	var err error
	cfg, err = config.Setup()
	if err != nil {
		log.Fatalf("Failed to setup config: %v", err)
	}

	httpClient = httpclient.NewHttpClient(&httpclient.Config{
		Timeout:          10 * time.Second,
		MaxRetries:       3,
		RetryWaitTime:    1 * time.Second,
		MaxRetryWaitTime: 10 * time.Second,
		DebugMode:        false,
	})

	ltiService = lti.NewService(cfg, httpClient)
	canvasService = canvas.NewService(cfg, httpClient)
}
