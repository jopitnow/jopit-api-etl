package utils

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jopitnow/go-jopit-toolkit/goutils/apierrors"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
)

// --- Canonical categories (fixed list in your platform) ---
var canonicalCategories = map[string]string{
	"10": "Home Decor",
	"11": "Furniture",
	"12": "Kitchen",
	"13": "Electronics",
}

func ValidateHexID(ids []string) apierrors.ApiError {

	// Regular expression to check if a string is a valid hex

	regex := regexp.MustCompile("^[a-fA-F0-9]{24}$")

	for _, id := range ids {
		val := regex.MatchString(id)
		if !val {
			return apierrors.NewApiError("one or more of the provided ids are not a valid hex string", "bad_request", 400, apierrors.CauseList{})
		}
	}

	return nil
}

func ExtractFromCSV(fileHeader *multipart.FileHeader) ([]map[string]string, apierrors.ApiError) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, apierrors.NewApiError("error opening csv file", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, apierrors.NewApiError("error reading content of csv file", "internal_server_error", http.StatusInternalServerError, apierrors.CauseList{})
	}

	headers := rows[0]
	var records []map[string]string
	for _, row := range rows[1:] {
		rec := make(map[string]string)
		for i, val := range row {
			rec[headers[i]] = val
		}
		records = append(records, rec)
	}
	return records, nil
}

// --- Transform with category mapping ---
func Transform(records []map[string]string, config models.CompanyLayout, userID string) (string, []models.Item) {
	var items []models.Item

	batchID := generateBatchID(6)

	for _, rec := range records {
		externalCat := rec[config.CategoryMap["category_name"]]
		mappedID := config.CategoryMap[externalCat]
		mappedName := canonicalCategories[mappedID]

		item := models.Item{
			ID:          rec[config.CategoryMap["id"]],
			ShopID:      config.ShopID,
			UserID:      userID,
			Name:        rec[config.CategoryMap["name"]],
			Description: rec[config.CategoryMap["description"]],
			Status:      "active",
			BatchID:     batchID,
			Fragile:     strings.ToLower(rec[config.CategoryMap["fragile"]]) == "true",
			Category: models.Category{
				ID:   mappedID,
				Name: mappedName,
			},
			Dimensions: models.Dimensions{
				Weight: parseInt(rec[config.CategoryMap["weight"]]),
				Length: parseInt(rec[config.CategoryMap["length"]]),
				Height: parseInt(rec[config.CategoryMap["height"]]),
				Width:  parseInt(rec[config.CategoryMap["width"]]),
			},
		}
		items = append(items, item)
	}
	return batchID, items
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func generateBatchID(length int) string {
	rand.Seed(int64(time.Now().UnixNano())) // Seed once here

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	batchID := make([]byte, length)
	for i := range batchID {
		batchID[i] = charset[rand.Intn(len(charset))]
	}
	return string(batchID)
}
