package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/repositories"
)

type CompanyLayoutService interface {
	Get(ctx context.Context, companyLayoutID string) (models.CompanyLayout, apierrors.ApiError)
	GetByShopID(ctx context.Context, shopID string) (models.CompanyLayout, apierrors.ApiError)
	GetAllCompanyLayout(ctx context.Context) ([]models.CompanyLayout, apierrors.ApiError)
	Create(ctx context.Context, input dto.CompanyLayoutRequest) apierrors.ApiError
	Update(context.Context, dto.CompanyLayoutRequest) apierrors.ApiError
	Delete(ctx context.Context) apierrors.ApiError
}

type companyLayout struct {
	repository  repositories.CompanyLayoutRepository
	shopsClient clients.ShopClient
}

func NewCompanyLayoutService(repository repositories.CompanyLayoutRepository, shopsClient clients.ShopClient) CompanyLayoutService {
	return &companyLayout{
		repository:  repository,
		shopsClient: shopsClient,
	}
}

func (s *companyLayout) Get(ctx context.Context, companyLayoutID string) (models.CompanyLayout, apierrors.ApiError) {
	companyLayout, err := s.repository.Get(ctx, companyLayoutID)
	if err != nil {
		return models.CompanyLayout{}, err
	}

	return companyLayout, nil
}

func (s *companyLayout) GetByShopID(ctx context.Context, shopID string) (models.CompanyLayout, apierrors.ApiError) {
	companyLayout, err := s.repository.Get(ctx, shopID)
	if err != nil {
		return models.CompanyLayout{}, err
	}

	return companyLayout, nil
}

func (s *companyLayout) GetAllCompanyLayout(ctx context.Context) ([]models.CompanyLayout, apierrors.ApiError) {
	companyLayouts, err := s.repository.GetAllCompanyLayout(ctx)
	if err != nil {
		return nil, err
	}

	if len(companyLayouts) <= 0 {
		return []models.CompanyLayout{}, apierrors.NewApiError("no companyLayout found", "companyLayout should never be nil, please contact and administrator", http.StatusInternalServerError, apierrors.CauseList{})
	}

	return companyLayouts, nil
}

func (s *companyLayout) Create(ctx context.Context, input dto.CompanyLayoutRequest) apierrors.ApiError {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return err
	}

	companyLayout := input.ToModel()
	if err != nil {
		return err
	}

	companyLayout.ShopID = shop.ID
	companyLayout.UserID = fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	_, err = s.repository.Create(ctx, companyLayout)
	if err != nil {
		return err
	}

	return nil
}

func (s *companyLayout) Update(ctx context.Context, input dto.CompanyLayoutRequest) apierrors.ApiError {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return err
	}

	companyLayout := input.ToModel()
	if err != nil {
		return err
	}

	companyLayout.ShopID = shop.ID
	companyLayout.UserID = fmt.Sprint(ctx.Value(goauth.FirebaseUserID))

	_, err = s.repository.Update(ctx, companyLayout)
	if err != nil {
		return err
	}

	return nil
}

func (s *companyLayout) Delete(ctx context.Context) apierrors.ApiError {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return err
	}

	companyLayout, err := s.repository.GetByShopID(ctx, shop.ID)
	if err != nil {
		return err
	}

	_, err = s.repository.Delete(ctx, companyLayout.ID)
	if err != nil {
		return err
	}

	return nil
}
