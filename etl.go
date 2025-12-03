package main

/*
import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Canonical categories (fixed list in your platform) ---
var canonicalCategories = map[string]string{
	"10": "Home Decor",
	"11": "Furniture",
	"12": "Kitchen",
	"13": "Electronics",
}

// --- Company config with category mapping ---
type CompanyLayout struct {
	UID         string            `bson:"uid"`
	ShopID      string            `bson:"shop_id"`
	Endpoint    string            `bson:"endpoint"`
	APIKey      string            `bson:"api_key"`
	FieldMap    map[string]string `bson:"field_map"`
	CategoryMap map[string]string `bson:"category_map"` // external → canonical ID
}

// --- Mongo init ---
func mongoClient() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func extractFromAPI(cfg RequestConfig) ([]map[string]string, error) {
	// Build request
	req, err := http.NewRequest(cfg.Method, cfg.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	// Add query params
	q := req.URL.Query()
	for k, v := range cfg.QueryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Add body if present
	if cfg.Method == "POST" || cfg.Method == "PUT" {
		bodyBytes, _ := json.Marshal(cfg.Body)
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var records []map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}

// --- Extract from CSV (accepts uploaded file directly) ---
func extractFromCSV(fileHeader *multipart.FileHeader) ([]map[string]string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
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
func transform(records []map[string]string, config CompanyLayout, userID string) []models.Item {
	var items []models.Item
	for _, rec := range records {
		externalCat := rec[config.FieldMap["category_name"]]
		mappedID := config.CategoryMap[externalCat]
		mappedName := canonicalCategories[mappedID]

		item := models.Item{
			ID:          rec[config.FieldMap["id"]],
			ShopID:      config.ShopID,
			UserID:      userID,
			Name:        rec[config.FieldMap["name"]],
			Description: rec[config.FieldMap["description"]],
			Status:      "active",
			Fragile:     strings.ToLower(rec[config.FieldMap["fragile"]]) == "true",
			Category: models.Category{
				ID:   mappedID,
				Name: mappedName,
			},
			Dimensions: models.Dimensions{
				Weight: parseInt(rec[config.FieldMap["weight"]]),
				Length: parseInt(rec[config.FieldMap["length"]]),
				Height: parseInt(rec[config.FieldMap["height"]]),
				Width:  parseInt(rec[config.FieldMap["width"]]),
			},
		}
		items = append(items, item)
	}
	return items
}

// --- Load with upsert ---
func load(items []models.Item) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mongoClient()
	defer client.Disconnect(ctx)

	collection := client.Database("yourdb").Collection("items")

	for _, item := range items {
		filter := bson.M{"_id": item.ID, "shop_id": item.ShopID}
		update := bson.M{"$set": item}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Println("Upsert failed for item:", item.ID, err)
		}
	}
	return nil
}

// --- Helpers ---
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// --- Gin API setup ---
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Register company API config
	r.POST("/register", func(c *gin.Context) {
		var cfg CompanyLayout
		if err := c.BindJSON(&cfg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		client := mongoClient()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := client.Database("yourdb").Collection("companies_config").InsertOne(ctx, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Company registered"})
	})

	// Trigger fetch from company API
	r.POST("/etl/fetch", func(c *gin.Context) {
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		token, err := firebaseAuth.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		uid := token.UID

		client := mongoClient()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var cfg CompanyLayout
		err = client.Database("yourdb").Collection("companies_config").FindOne(ctx, bson.M{"uid": uid}).Decode(&cfg)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Company not registered"})
			return
		}

		records, err := extractFromAPI(cfg.Endpoint, cfg.APIKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Extract failed"})
			return
		}

		items := transform(records, cfg, uid)
		if err := load(items); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Load failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ETL completed"})
	})

	// CSV upload
	r.POST("/etl/csv", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file required"})
			return
		}

		// Pass the fileHeader directly to extractFromCSV
		records, err := extractFromCSV(fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Extract failed"})
			return
		}

		// Lookup company config (based on auth UID, etc.)
		var cfg CompanyLayout
		// ... fetch from MongoDB ...

		items := transform(records, cfg, "userID-from-auth")
		if err := load(items); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Load failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ETL completed"})
	})

	return r
}
*/
