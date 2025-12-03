package dependencies

import (
	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/jopitnow/jopit-api-etl/src/main/api/dependencies"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
	"github.com/jopitnow/jopit-api-etl/src/tests/internal/api/platform/storage"
	"github.com/tryvium-travels/memongo"
)

type DependencyManagerMock struct {
	*gonosql.Data
}

func NewMockDependencyManager(server *memongo.Server) DependencyManagerMock {
	db := storage.OpenNoSQLMock(server)
	if db.Error != nil {
		panic(db.Error)
	}

	return DependencyManagerMock{
		db,
	}
}

func (m DependencyManagerMock) ItemsRepository() repositories.ItemsRepository {
	return repositories.NewItemsRepository(m.NewCollection(dependencies.KvsItemsCollection))
}

func (m DependencyManagerMock) CategoriesRepository() repositories.CategoriesRepository {
	return repositories.NewCategoriesRepository(m.NewCollection(dependencies.KvsCategoriesCollection))
}
