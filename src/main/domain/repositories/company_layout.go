package repositories

import (
	"context"
	"errors"
	"net/http"

	"go.opentelemetry.io/otel"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	CompanyLayoutDatabaseError = "[%s] Error in DB"
)

var tracerRepoCompanyLayout = otel.Tracer("companyLayout-repo") // Tracer for this package

type CompanyLayoutRepository interface {
	Get(ctx context.Context, companyLayoutID string) (models.CompanyLayout, apierrors.ApiError)
	GetByShopID(ctx context.Context, shopID string) (models.CompanyLayout, apierrors.ApiError)
	GetAllCompanyLayout(ctx context.Context) ([]models.CompanyLayout, apierrors.ApiError)
	Create(ctx context.Context, input models.CompanyLayout) (interface{}, apierrors.ApiError)
	Update(ctx context.Context, input models.CompanyLayout) (int64, apierrors.ApiError)
	Delete(ctx context.Context, companyLayoutID string) (int64, apierrors.ApiError)
}

type companyLayoutRepository struct {
	Collection *mongo.Collection
}

func NewCompanyLayoutRepository(Collection *mongo.Collection) CompanyLayoutRepository {
	return &companyLayoutRepository{Collection: Collection}
}

func (storage *companyLayoutRepository) Get(ctx context.Context, companyLayoutID string) (models.CompanyLayout, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "Get")
	defer span.End()

	var companyLayout models.CompanyLayout

	result, err := gonosql.Get(ctx, storage.Collection, companyLayoutID)
	if err != nil {
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error", "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	if result.Err() != nil { // coverage-ignore
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error: "+result.Err().Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	err = result.Decode(&companyLayout)
	if err != nil {
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	return companyLayout, nil
}

func (storage *companyLayoutRepository) GetByShopID(ctx context.Context, shopID string) (models.CompanyLayout, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "Get")
	defer span.End()

	var companyLayout models.CompanyLayout

	filter := bson.M{"shop_id": shopID}

	result := storage.Collection.FindOne(ctx, filter)

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error", "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	if result.Err() != nil { // coverage-ignore
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error: "+result.Err().Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	err := result.Decode(&companyLayout)
	if err != nil {
		return models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Get() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	return companyLayout, nil
}

func (storage *companyLayoutRepository) GetAllCompanyLayout(ctx context.Context) ([]models.CompanyLayout, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "GetAllCompanyLayout")
	defer span.End()

	var companyLayout []models.CompanyLayout

	filter := bson.M{}

	cursor, err := storage.Collection.Find(ctx, filter)
	if err != nil {
		return []models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("GetAllCompanyLayout() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if err = cursor.All(ctx, &companyLayout); err != nil {
		return []models.CompanyLayout{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("GetAllCompanyLayout() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	return companyLayout, nil
}

func (storage *companyLayoutRepository) Create(ctx context.Context, input models.CompanyLayout) (interface{}, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "Create")
	defer span.End()

	result, err := gonosql.InsertOne(ctx, storage.Collection, input)
	if err != nil {
		return nil, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Create() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if result.InsertedID == nil || result.InsertedID == "" { // coverage-ignore
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Create() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	return result.InsertedID, nil
}

func (storage *companyLayoutRepository) Update(ctx context.Context, input models.CompanyLayout) (int64, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "Update")
	defer span.End()

	primitiveID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Update() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	update := bson.M{
		"$set": bson.M{
			"item_map":     input.ItemMap,
			"category_map": input.CategoryMap,
		},
	}

	result, err := storage.Collection.UpdateOne(ctx, bson.M{"_id": primitiveID}, update)
	if err != nil {
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Update() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if result.MatchedCount == 0 {
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Update() error", "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	return result.ModifiedCount, nil
}

func (storage *companyLayoutRepository) Delete(ctx context.Context, companyLayoutID string) (int64, apierrors.ApiError) {

	ctx, span := tracerRepoCompanyLayout.Start(ctx, "Delete")
	defer span.End()

	result, err := gonosql.Delete(ctx, storage.Collection, companyLayoutID)
	if err != nil {
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Delete() error: "+err.Error(), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{}))
	}

	if result.DeletedCount == 0 {
		return -1, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError("Delete() error", "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	return result.DeletedCount, nil
}
