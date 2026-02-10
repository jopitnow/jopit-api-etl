package dependencies

import (
	"github.com/jopitnow/jopit-api-etl/src/main/api/platform/storage"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
)

const (
	KvsCompanyLayoutCollection = "company-layout"
	KvsMercadoLibreCredentials = "mercadolibre-credentials"
)

type DependencyManager struct {
	*gonosql.Data
}

func NewDependencyManager() DependencyManager {
	db := storage.NewNoSQL()
	if db.Error != nil {
		panic(db.Error)
	}
	return DependencyManager{
		db,
	}
}

func (m DependencyManager) CompanyLayoutRepository() repositories.CompanyLayoutRepository {
	return repositories.NewCompanyLayoutRepository(m.NewCollection(KvsCompanyLayoutCollection))
}

func (m DependencyManager) MercadoLibreCredentialsRepository() repositories.MercadoLibreCredentialsRepository {
	return repositories.NewMercadoLibreCredentialsRepository(m.NewCollection(KvsMercadoLibreCredentials))
}
