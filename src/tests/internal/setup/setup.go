package setup

import (
	"fmt"

	"github.com/jopitnow/jopit-api-etl/src/tests/internal/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/api/platform/storage"
	"github.com/tryvium-travels/memongo"
)

var server *memongo.Server

func BeforeMemongoTestCase() dependencies.Dependencies {
	var err error
	server, err = memongo.Start("6.0.19")
	if err != nil {
		fmt.Println("Error starting on memory MongoDB server.", err)
	}

	deps, err := dependencies.BuildMockDependencies(server)
	if err != nil {
		fmt.Println("Error creating repository.", err)
	}

	return deps
}

func AfterMemongoTestCase() {
	storage.CloseNoSQLMock()
	server.Stop()
}
