package dependencies

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/handlers"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"
)

type Dependencies interface {
	CompanyLayoutRepository() repositories.CompanyLayoutRepository
}

func GetDependencyManager() Dependencies {
	return NewDependencyManager()
}

func BuildDependencies() (HandlersStruct, error) {

	manager := GetDependencyManager()

	caompanyLayoutRepository := manager.CompanyLayoutRepository()

	// External Clients
	fetchApiClient := clients.FetchApiClientInstance
	itemsClient := clients.ItemsClientInstance
	shopsClient := clients.ShopsClientInstance

	// Services
	etlService := services.NewEtlService(fetchApiClient, itemsClient, shopsClient)
	companyLayoutService := services.NewCompanyLayoutService(caompanyLayoutRepository, shopsClient)

	// Handlers
	etlHandler := handlers.NewEtlsHandler(etlService)
	companyLayoutHandler := handlers.NewCompanyLayoutHandler(companyLayoutService)

	return HandlersStruct{
		Etl:           etlHandler,
		CompanyLayout: companyLayoutHandler,
	}, nil
}

type HandlersStruct struct {
	Etl           handlers.EtlHandler
	CompanyLayout handlers.CompanyLayoutHandler
}
