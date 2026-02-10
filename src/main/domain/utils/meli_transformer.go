package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/jopitnow/jopit-api-etl/src/main/domain/models"
	"github.com/jopitnow/jopit-api-etl/src/main/domain/models/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransformMeliItemToJopitItem converts a MercadoLibre item to Jopit format
func TransformMeliItemToJopitItem(
	meliItem dto.MeliItemResponse,
	shopID string,
	userID string,
	batchID string,
	sizeChart *dto.MeliSizeChartResponse,
) models.Item {

	item := models.Item{
		ID:          primitive.NewObjectID().Hex(),
		ShopID:      shopID,
		UserID:      userID,
		Name:        meliItem.Title,
		Description: extractDescription(meliItem),
		Status:      mapStatus(meliItem.Status),
		Category:    mapCategory(meliItem.CategoryID, meliItem.DomainID),
		Delivery:    mapDelivery(meliItem.Shipping),
		Attributes:  extractAttributes(meliItem.Attributes, meliItem.Condition, meliItem.SaleTerms),
		Variants:    mapVariants(meliItem.Variations, meliItem.Pictures),
		Price:       mapPrice(meliItem, shopID),
		Source: &models.Source{
			SourceType:        "meli",
			ExternalID:        meliItem.ID,
			ExternalSKU:       extractExternalSKU(meliItem.Variations),
			BatchID:           batchID,
			ImportedAt:        time.Now(),
			EtlVersion:        "1.0.0",
			TransformMetadata: mapTransformMetadata(meliItem),
		},
	}

	// Map size guide if available
	if sizeChart != nil {
		item.SizeGuide = mapSizeGuide(*sizeChart)
		// Store SIZE_GRID_ID in SizeGuide
		sizeGridID := ExtractSizeChartID(meliItem.Attributes)
		if sizeGridID != "" {
			item.SizeGuide.ExternalSizeGridID = sizeGridID
		}
	}

	item.ValidateEmptySlices()
	return item
}

// ExtractAttributeValue extracts a specific attribute value by ID
func ExtractAttributeValue(attributes []dto.MeliAttribute, attributeID string) string {
	for _, attr := range attributes {
		if attr.ID == attributeID {
			if attr.ValueName != "" {
				return attr.ValueName
			}
			if len(attr.Values) > 0 {
				return attr.Values[0].Name
			}
		}
	}
	return ""
}

// extractDescription builds description from attributes
func extractDescription(meliItem dto.MeliItemResponse) string {
	brand := ExtractAttributeValue(meliItem.Attributes, "BRAND")
	if brand != "" {
		return fmt.Sprintf("%s - %s", brand, meliItem.Title)
	}
	return meliItem.Title
}

// mapStatus converts MercadoLibre status to Jopit status
func mapStatus(meliStatus string) string {
	switch meliStatus {
	case "active":
		return "active"
	case "paused", "inactive":
		return "inactive"
	case "closed":
		return "inactive"
	default:
		return "active"
	}
}

// mapCategory maps MercadoLibre category to Jopit category
func mapCategory(categoryID string, domainID string) models.ItemCategory {
	// For now, create a generic category based on domain
	// TODO: Implement proper category mapping from configuration
	categoryName := "Ropa"
	if strings.Contains(domainID, "SHOE") || strings.Contains(domainID, "FOOTWEAR") {
		categoryName = "Calzado"
	} else if strings.Contains(domainID, "ACCESSORIE") {
		categoryName = "Accesorios"
	}

	cat := models.ItemCategory{
		ID:   primitive.NewObjectID().Hex(),
		Name: categoryName,
	}
	cat.SetIDs()
	return cat
}

// mapDelivery converts shipping info to delivery format
func mapDelivery(shipping dto.MeliShipping) models.Delivery {
	// Default dimensions if not provided
	return models.Delivery{
		Fragile: false,
		Dimensions: models.Dimensions{
			Weight: 500, // 500g default
			Length: 30,  // 30cm default
			Height: 5,   // 5cm default
			Width:  25,  // 25cm default
		},
	}
}

// extractAttributes extracts product attributes as struct and MercadoLibre-specific attributes
func extractAttributes(attributes []dto.MeliAttribute, condition string, saleTerms []interface{}) models.Attributes {
	gender := ExtractAttributeValue(attributes, "GENDER")
	composition := ExtractAttributeValue(attributes, "COMPOSITION")

	// Map condition
	jopitCondition := "new"
	if condition == "used" {
		jopitCondition = "pre-owned"
	}

	// Extract MercadoLibre attributes (exclude specific ones)
	excludedAttrs := map[string]bool{
		"GIFTABLE":                 true,
		"IS_EMERGING_BRAND":        true,
		"IS_HIGHLIGHT_BRAND":       true,
		"IS_SUITABLE_FOR_PREGNACY": true,
		"IS_TOM_BRAND":             true,
		"SIZE_GRID_ID":             true,
		"WITH_RECYCLED_MATERIALS":  true,
	}

	meliAttrs := make([]models.MeliAttribute, 0)

	// Add MercadoLibre attributes
	for _, attr := range attributes {
		// Skip excluded attributes
		if excludedAttrs[attr.ID] {
			continue
		}

		// Skip attributes with empty value_name
		if attr.ValueName == "" {
			continue
		}

		// Skip attributes with invalid value_id (-1 indicates unset/invalid)
		if attr.ValueID != nil && *attr.ValueID == "-1" {
			continue
		}

		valueID := ""
		if attr.ValueID != nil {
			valueID = *attr.ValueID
		}

		meliAttrs = append(meliAttrs, models.MeliAttribute{
			ID:        attr.ID,
			Name:      attr.Name,
			ValueID:   valueID,
			ValueName: attr.ValueName,
		})
	}

	// Add sale terms as MeliAttributes with "sale_term_" prefix
	for _, term := range saleTerms {
		if termMap, ok := term.(map[string]interface{}); ok {
			if id, ok := termMap["id"].(string); ok {
				if valueName, ok := termMap["value_name"].(string); ok {
					valueID := ""
					if vid, ok := termMap["value_id"].(string); ok {
						valueID = vid
					}

					name := ""
					if n, ok := termMap["name"].(string); ok {
						name = n
					}

					meliAttrs = append(meliAttrs, models.MeliAttribute{
						ID:        fmt.Sprintf("sale_term_%s", strings.ToLower(id)),
						Name:      name,
						ValueID:   valueID,
						ValueName: valueName,
					})
				}
			}
		}
	}

	return models.Attributes{
		Condition:      jopitCondition,
		Gender:         strings.ToLower(gender),
		Composition:    composition,
		MeliAttributes: meliAttrs,
	}
}

// mapVariants converts MercadoLibre variations to Jopit variants
func mapVariants(variations []dto.MeliVariation, pictures []dto.MeliPicture) []models.Variant {
	if len(variations) == 0 {
		// No variations, create a single variant with all images
		return []models.Variant{
			{
				ColorID:   "default",
				ColorName: "Default",
				ColorHex:  "#000000",
				IsMain:    true,
				Images:    extractImageURLs(pictures),
				SizeStock: []models.SizeStock{},
			},
		}
	}

	// Group variations by color
	colorMap := make(map[string]*models.Variant)

	for _, variation := range variations {
		colorID, colorName := extractColorFromVariation(variation)
		sizeLabel := extractSizeFromVariation(variation)

		// Get or create variant for this color
		variant, exists := colorMap[colorID]
		if !exists {
			variant = &models.Variant{
				ColorID:   colorID,
				ColorName: colorName,
				ColorHex:  "#000000", // TODO: Map color names to hex
				IsMain:    len(colorMap) == 0,
				Images:    mapVariationImages(variation, pictures),
				SizeStock: []models.SizeStock{},
			}
			colorMap[colorID] = variant
		}

		// Add size stock
		if sizeLabel != "" {
			variant.SizeStock = append(variant.SizeStock, models.SizeStock{
				SizeLabel: sizeLabel,
				Stock:     variation.AvailableQuantity,
				SKU:       variation.UserProductID,
			})
		}
	}

	// Convert map to slice
	variants := make([]models.Variant, 0, len(colorMap))
	for _, variant := range colorMap {
		variants = append(variants, *variant)
	}

	return variants
}

// extractColorFromVariation extracts color info from variation attributes
func extractColorFromVariation(variation dto.MeliVariation) (string, string) {
	for _, attr := range variation.AttributeCombinations {
		if attr.ID == "COLOR" || attr.ID == "MAIN_COLOR" {
			if attr.ValueID != nil && *attr.ValueID != "" {
				return *attr.ValueID, attr.ValueName
			}
			return attr.ValueName, attr.ValueName
		}
	}
	return "default", "Default"
}

// extractSizeFromVariation extracts size label from variation
func extractSizeFromVariation(variation dto.MeliVariation) string {
	for _, attr := range variation.AttributeCombinations {
		if attr.ID == "SIZE" {
			return attr.ValueName
		}
	}
	return ""
}

// mapVariationImages maps variation pictures to image URLs
func mapVariationImages(variation dto.MeliVariation, allPictures []dto.MeliPicture) []models.Image {
	images := []models.Image{}

	// Map picture IDs to URLs
	pictureMap := make(map[string]string)
	for _, pic := range allPictures {
		pictureMap[pic.ID] = pic.SecureURL
	}

	// Get images for this variation
	for _, picID := range variation.PictureIDs {
		if url, exists := pictureMap[picID]; exists {
			images = append(images, models.Image(url))
		}
	}

	// If no specific images, use all pictures
	if len(images) == 0 {
		images = extractImageURLs(allPictures)
	}

	return images
}

// extractImageURLs extracts all image URLs from pictures
func extractImageURLs(pictures []dto.MeliPicture) []models.Image {
	urls := make([]models.Image, 0, len(pictures))
	for _, pic := range pictures {
		urls = append(urls, models.Image(pic.SecureURL))
	}
	return urls
}

// mapSizeGuide converts MercadoLibre size chart to Jopit format
func mapSizeGuide(sizeChart dto.MeliSizeChartResponse) *models.SizeGuide {
	sizes := make([]models.Size, 0, len(sizeChart.Rows))

	for _, row := range sizeChart.Rows {
		size := models.Size{}

		for _, attr := range row.Attributes {
			switch attr.ID {
			case "SIZE":
				if len(attr.Values) > 0 {
					size.SizeEquivalence = attr.Values[0].Name
				}
			case "CHEST_CIRCUMFERENCE_FROM":
				if len(attr.Values) > 0 && attr.Values[0].Struct != nil {
					if ms := attr.Values[0].Struct; ms != nil {
						size.ChestCircumference = int(ms.Number)
					}
				}
			case "WAIST_CIRCUMFERENCE_FROM":
				if len(attr.Values) > 0 && attr.Values[0].Struct != nil {
					if ms := attr.Values[0].Struct; ms != nil {
						size.WaistCircumference = int(ms.Number)
					}
				}
			case "HIP_CIRCUMFERENCE_FROM":
				if len(attr.Values) > 0 && attr.Values[0].Struct != nil {
					if ms := attr.Values[0].Struct; ms != nil {
						size.HipCircumference = int(ms.Number)
					}
				}
			}
		}

		if size.SizeEquivalence != "" {
			sizes = append(sizes, size)
		}
	}

	// Determine body part from measure type
	bodyPart := "upper" // default
	if sizeChart.MeasureType == "BODY_MEASURE" {
		bodyPart = "upper"
	}

	return &models.SizeGuide{
		Type:              "standard",
		BodyPart:          bodyPart,
		HasMeasurements:   len(sizes) > 0,
		IsOneSize:         len(sizes) == 1,
		MeasurementSource: "mercadolibre",
		Sizes:             sizes,
	}
}

// extractExternalSKU extracts the main SKU from variations
func extractExternalSKU(variations []dto.MeliVariation) string {
	if len(variations) > 0 && variations[0].UserProductID != "" {
		return variations[0].UserProductID
	}
	return ""
}

// mapTransformMetadata creates metadata about the ETL transformation
func mapTransformMetadata(meliItem dto.MeliItemResponse) map[string]string {
	metadata := make(map[string]string)

	metadata["original_domain_id"] = meliItem.DomainID
	metadata["original_category_id"] = meliItem.CategoryID
	metadata["original_status"] = meliItem.Status
	metadata["original_condition"] = meliItem.Condition
	metadata["variations_count"] = fmt.Sprintf("%d", len(meliItem.Variations))
	metadata["pictures_count"] = fmt.Sprintf("%d", len(meliItem.Pictures))
	metadata["attributes_count"] = fmt.Sprintf("%d", len(meliItem.Attributes))

	return metadata
}

// mapPrice extracts price information from MercadoLibre item
func mapPrice(meliItem dto.MeliItemResponse, shopID string) models.Price {
	return models.Price{
		ShopID:   shopID,
		Amount:   meliItem.Price,
		Currency: mapCurrency(meliItem.CurrencyID),
	}
}

// mapCurrency creates a Currency struct from MercadoLibre currency ID
func mapCurrency(currencyID string) models.Currency {
	// Map MercadoLibre currency codes to Jopit currency info
	currencyMap := map[string]models.Currency{
		"ARS": {
			ID:               "ARS",
			Symbol:           "$",
			DecimalDivider:   ",",
			ThousandsDivider: ".",
		},
		// "USD": {
		// 	ID:               "USD",
		// 	Symbol:           "$",
		// 	DecimalDivider:   ".",
		// 	ThousandsDivider: ",",
		// },
		// "BRL": {
		// 	ID:               "BRL",
		// 	Symbol:           "R$",
		// 	DecimalDivider:   ",",
		// 	ThousandsDivider: ".",
		// },
		// "MXN": {
		// 	ID:               "MXN",
		// 	Symbol:           "$",
		// 	DecimalDivider:   ".",
		// 	ThousandsDivider: ",",
		// },
		// "CLP": {
		// 	ID:               "CLP",
		// 	Symbol:           "$",
		// 	DecimalDivider:   ",",
		// 	ThousandsDivider: ".",
		// },
		// "UYU": {
		// 	ID:               "UYU",
		// 	Symbol:           "$U",
		// 	DecimalDivider:   ",",
		// 	ThousandsDivider: ".",
		// },
		// "PEN": {
		// 	ID:               "PEN",
		// 	Symbol:           "S/",
		// 	DecimalDivider:   ".",
		// 	ThousandsDivider: ",",
		// },
		// "COP": {
		// 	ID:               "COP",
		// 	Symbol:           "$",
		// 	DecimalDivider:   ",",
		// 	ThousandsDivider: ".",
		// },
		// "EUR": {
		// 	ID:               "EUR",
		// 	Symbol:           "€",
		// 	DecimalDivider:   ",",
		// 	ThousandsDivider: ".",
		// },
	}

	if currency, exists := currencyMap[currencyID]; exists {
		return currency
	}

	// Default fallback for unknown currencies
	return models.Currency{
		ID:               currencyID,
		Symbol:           currencyID,
		DecimalDivider:   ".",
		ThousandsDivider: ",",
	}
}

// ExtractSizeChartID extracts SIZE_GRID_ID from item attributes
func ExtractSizeChartID(attributes []dto.MeliAttribute) string {
	for _, attr := range attributes {
		if attr.ID == "SIZE_GRID_ID" {
			if attr.ValueName != "" {
				return attr.ValueName
			}
			if attr.ValueID != nil {
				return *attr.ValueID
			}
		}
	}
	return ""
}
