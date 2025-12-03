package dependencies

import (
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
	"github.com/tryvium-travels/memongo"
)

type DependencyMock interface {
	ItemsRepository() repositories.ItemsRepository
	CategoriesRepository() repositories.CategoriesRepository
}

func GetDependencyManagerMock(server *memongo.Server) DependencyMock {
	return NewMockDependencyManager(server)
}

func BuildMockDependencies(server *memongo.Server) (Dependencies, error) {
	manager := GetDependencyManagerMock(server)

	itemsRepository := manager.ItemsRepository()
	categoriesRepository := manager.CategoriesRepository()

	return Dependencies{
		ItemsRepository:      itemsRepository,
		CategoriesRepository: categoriesRepository,
	}, nil
}

type Dependencies struct {
	ItemsRepository      repositories.ItemsRepository
	CategoriesRepository repositories.CategoriesRepository
}
