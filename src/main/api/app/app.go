package app

import (
	"fmt"
	"time"

	"github.com/jopitnow/jopit-api-etl/src/main/api/config"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"

	"github.com/gin-gonic/gin"
	"github.com/jopitnow/go-jopit-toolkit/gingonic/handlers"
	"github.com/jopitnow/go-jopit-toolkit/goutils/logger"
)

func Start() {
	handler, err := dependencies.BuildDependencies()
	if err != nil {
		fmt.Printf("Error Build Dependencies")
		waitAndPanic(err)

	}

	router := ConfigureRouter()
	RouterMapper(router, handler)

	if errRouter := router.Run(config.ConfMap.APIRestServerPort); errRouter != nil {
		logger.Errorf("Error starting router", errRouter)
		waitAndPanic(errRouter)
	}
}

func ConfigureRouter() *gin.Engine {
	logger.InitLog(config.ConfMap.LoggingPath, config.ConfMap.LoggingFile, config.ConfMap.LoggingLevel)
	return handlers.DefaultJopitRouter()

}

func waitAndPanic(err error) {
	time.Sleep(2 * time.Second) // needs one second to send the log
	panic(err)
}
