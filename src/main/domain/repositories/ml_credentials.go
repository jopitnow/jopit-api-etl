package repositories

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/mgo.v2/bson"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MercadoLibreDatabaseError = "[%s] Error in DB"
	dbErr                     = MercadoLibreDatabaseError
)

var tracerMeliRepo = otel.Tracer("meli-credentials-repo")

type MercadoLibreCredentialsRepository interface {
	GetCredentialsByShopID(ctx context.Context, shopID string) (models.MercadoLibreCredential, apierrors.ApiError)
	GetCredentialsByUserID(ctx context.Context, userID string) (models.MercadoLibreCredential, apierrors.ApiError)
	CreateCredentials(ctx context.Context, credentials models.MercadoLibreCredential) apierrors.ApiError
	UpdateCredentials(ctx context.Context, credentials models.MercadoLibreCredential) apierrors.ApiError
	DeleteCredentials(ctx context.Context, userID string) apierrors.ApiError
}

type mercadoLibreCredentialsRepository struct {
	Collection *mongo.Collection
}

func NewMercadoLibreCredentialsRepository(collection *mongo.Collection) MercadoLibreCredentialsRepository {
	return &mercadoLibreCredentialsRepository{
		Collection: collection,
	}
}

func (r *mercadoLibreCredentialsRepository) GetCredentialsByShopID(ctx context.Context, shopID string) (models.MercadoLibreCredential, apierrors.ApiError) {
	ctx, span := tracerMeliRepo.Start(ctx, "GetCredentialsByShopID")
	defer span.End()

	filter := bson.M{"shop_id": shopID}

	return r.getCredentials(ctx, filter, span)
}

func (r *mercadoLibreCredentialsRepository) GetCredentialsByUserID(ctx context.Context, userID string) (models.MercadoLibreCredential, apierrors.ApiError) {
	ctx, span := tracerMeliRepo.Start(ctx, "GetCredentialsByUserID")
	defer span.End()

	filter := bson.M{"user_id": userID}

	return r.getCredentials(ctx, filter, span)
}

func (r *mercadoLibreCredentialsRepository) CreateCredentials(ctx context.Context, credentials models.MercadoLibreCredential) apierrors.ApiError {
	ctx, span := tracerMeliRepo.Start(ctx, "CreateCredentials")
	defer span.End()

	result, err := gonosql.InsertOne(ctx, r.Collection, credentials)
	if err != nil {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "CreateCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	if result.InsertedID == nil || result.InsertedID == "" {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "CreateCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return nil
}

func (r *mercadoLibreCredentialsRepository) UpdateCredentials(ctx context.Context, credentials models.MercadoLibreCredential) apierrors.ApiError {
	ctx, span := tracerMeliRepo.Start(ctx, "UpdateCredentials")
	defer span.End()

	primitiveID, err := primitive.ObjectIDFromHex(credentials.ID)
	if err != nil {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "UpdateCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err.Error()}))
	}

	credentials.ID = ""

	update := bson.M{"$set": credentials}

	result, err := r.Collection.UpdateOne(ctx, bson.M{"_id": primitiveID}, update)
	if err != nil {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "UpdateCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err.Error()}))
	}

	if result.MatchedCount == 0 {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(dbErr, "UpdateCredentials"), "not_found", http.StatusNotFound, apierrors.CauseList{}))
	}

	return nil
}

func (r *mercadoLibreCredentialsRepository) DeleteCredentials(ctx context.Context, userID string) apierrors.ApiError {
	ctx, span := tracerMeliRepo.Start(ctx, "DeleteCredentials")
	defer span.End()

	result, err := r.Collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(dbErr, "DeleteCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	if result.DeletedCount == 0 {
		return apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(dbErr, "DeleteCredentials"), "not_found", http.StatusNotFound, apierrors.CauseList{err}))
	}

	return nil
}

func (r *mercadoLibreCredentialsRepository) getCredentials(ctx context.Context, filter bson.M, span trace.Span) (models.MercadoLibreCredential, apierrors.ApiError) {
	var model models.MercadoLibreCredential
	result := r.Collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return models.MercadoLibreCredential{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "GetCredentials"), "not_found", http.StatusNotFound, apierrors.CauseList{"no documents found"}))
	}

	if result.Err() != nil {
		return models.MercadoLibreCredential{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "GetCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{result.Err()}))
	}

	err := result.Decode(&model)
	if err != nil {
		return models.MercadoLibreCredential{}, apierrors.NewWrapAndTraceError(span, apierrors.NewApiError(fmt.Sprintf(MercadoLibreDatabaseError, "GetCredentials"), "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{err}))
	}

	return model, nil
}
