package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/clients"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/utils"

	"github.com/jopitnow/go-jopit-toolkit/goauth"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
)

type EtlService interface {
	LoadApi(ctx context.Context) (string, apierrors.ApiError)
	LoadCsv(ctx context.Context, file *multipart.FileHeader) (string, apierrors.ApiError)
	DeleteBatch(ctx context.Context, batchID string) apierrors.ApiError
}

type etlService struct {
	companyConfigService CompanyLayoutService
	itemsClient          clients.ItemsClient
	httpClient           clients.FetchApiClient
	shopsClient          clients.ShopClient
}

func NewEtlService(
	httpClient clients.FetchApiClient,
	itemsClient clients.ItemsClient,
	shopsClient clients.ShopClient,
) EtlService {
	return &etlService{
		httpClient:  httpClient,
		itemsClient: itemsClient,
		shopsClient: shopsClient,
	}
}

func (s *etlService) LoadApi(ctx context.Context) (string, apierrors.ApiError) {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return "", err
	}

	companyLayout, err := s.companyConfigService.GetByShopID(ctx, shop.ID)
	if err != nil {
		return "", err
	}

	//validate Layout ??

	response, err := s.httpClient.FetchAPI(ctx, companyLayout)
	if err != nil {
		return "", err
	}

	batchID, items := utils.Transform(response, companyLayout, fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	err = s.itemsClient.BulkCreateItems(ctx, items)
	if err != nil {
		return "", err
	}

	return batchID, nil
}

func (s *etlService) LoadCsv(ctx context.Context, file *multipart.FileHeader) (string, apierrors.ApiError) {

	shop, err := s.shopsClient.GetShopByUserID(ctx)
	if err != nil {
		return "", err
	}

	companyLayout, err := s.companyConfigService.GetByShopID(ctx, shop.ID)
	if err != nil {
		return "", err
	}

	//validate Layout ??

	data, err := utils.ExtractFromCSV(file)
	if err != nil {
		return "", err
	}

	batchID, items := utils.Transform(data, companyLayout, fmt.Sprint(ctx.Value(goauth.FirebaseAuthHeader)))

	err = s.itemsClient.BulkCreateItems(ctx, items)
	if err != nil {
		return "", err
	}

	return batchID, nil
}

func (s *etlService) DeleteBatch(ctx context.Context, batchID string) apierrors.ApiError {
	return s.itemsClient.BulkDeleteItems(ctx, batchID)
}
