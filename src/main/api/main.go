package main

import (
	"context"
	"log"

	toolkitpkg "github.com/jopitnow/go-jopit-toolkit/gingonic/handlers"

	"github.com/jopitnow/go-jopit-toolkit/telemetry"

	"github.com/jopitnow/jopit-api-etl/src/main/api/app"
	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
)

// @title Jopit API Items
// @version 1.0
// @description This is a jopit api items.
// @termsOfService http://swagger.io/terms/

// @contact.name Agustin Rabini
// @contact.url http://www.swagger.io/support
// @contact.email agustinrabini99@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	config.Load()

	shutdown, err := telemetry.InitTracerExporter(toolkitpkg.ApiName)
	if err != nil {
		log.Printf("failed to initialize tracer: %v", err)
	}
	// Ensure that all spans are flushed before the application exits.
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	shutdown, err = telemetry.InitLoggerExporter(toolkitpkg.ApiName)
	if err != nil {
		log.Printf("failed to initialize tracer: %v", err)
	}
	// Ensure that all spans are flushed before the application exits.
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	shutdown, err = telemetry.InitMeterExporter(toolkitpkg.ApiName)
	if err != nil {
		log.Printf("failed to initialize tracer: %v", err)
	}
	// Ensure that all spans are flushed before the application exits.
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	app.Start()
}
