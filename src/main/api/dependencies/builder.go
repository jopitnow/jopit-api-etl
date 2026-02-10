package dependencies

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/handlers"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/services"
)

type Dependencies interface {
	CompanyLayoutRepository() repositories.CompanyLayoutRepository
	MercadoLibreCredentialsRepository() repositories.MercadoLibreCredentialsRepository
}

func GetDependencyManager() Dependencies {
	return NewDependencyManager()
}

func BuildDependencies() (HandlersStruct, error) {

	manager := GetDependencyManager()

	caompanyLayoutRepository := manager.CompanyLayoutRepository()
	mercadoLibreCredentialsRepository := manager.MercadoLibreCredentialsRepository()

	// External Clients
	fetchApiClient := clients.FetchApiClientInstance
	itemsClient := clients.ItemsClientInstance
	shopsClient := clients.ShopsClientInstance
	mercadoLibreAuthClient := clients.MercadoLibreAuthClientInstance
	mercadoLibreClient := clients.MercadoLibreClientInstance

	// Services
	mercadoLibreCredentialsService := services.NewMercadoLibreCredentialsService(mercadoLibreCredentialsRepository, shopsClient, mercadoLibreAuthClient)
	mercadoLibreService := services.NewMercadoLibreService(mercadoLibreClient, mercadoLibreCredentialsService)
	etlService := services.NewEtlService(fetchApiClient, itemsClient, shopsClient)
	companyLayoutService := services.NewCompanyLayoutService(caompanyLayoutRepository, shopsClient)

	// Handlers
	etlHandler := handlers.NewEtlsHandler(etlService, mercadoLibreService)
	companyLayoutHandler := handlers.NewCompanyLayoutHandler(companyLayoutService)
	mercadoLibreCredentialsHandler := handlers.NewMercadoLibreCredentialsHandler(mercadoLibreCredentialsService)

	return HandlersStruct{
		Etl:                     etlHandler,
		CompanyLayout:           companyLayoutHandler,
		MercadoLibreCredentials: mercadoLibreCredentialsHandler,
	}, nil
}

type HandlersStruct struct {
	Etl                     handlers.EtlHandler
	CompanyLayout           handlers.CompanyLayoutHandler
	MercadoLibreCredentials handlers.MercadoLibreCredentialsHandler
}
