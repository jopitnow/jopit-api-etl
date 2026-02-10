# ETL Process Documentation

## Overview

The ETL (Extract, Transform, Load) process is the core functionality of the Jopit ETL service. It connects to external e-commerce platforms (currently MercadoLibre) to fetch product catalog data, transforms it into Jopit's standardized Item model, and loads it into the Jopit Items domain.

---

## The Three Phases

### 1. Extract Phase

**Purpose**: Retrieve raw product data from external sources with proper authentication.

**Implementation Location**: `src/main/domain/services/mercadolibre.go`

**Key Functions**:
- `GetUserItems()` - Fetches list of item IDs for a seller
- `GetItems()` - Batch fetches full details for multiple items
- `GetSizeChart()` - Retrieves size guide information

**Process**:
1. Authenticate using OAuth 2.0 tokens stored in MongoDB
2. Auto-refresh tokens if they're expiring within 1 hour (handled in `src/main/domain/clients/http.go`)
3. Call MercadoLibre API endpoints to fetch seller's product catalog
4. Return raw product data in MercadoLibre's format

**Output**: Array of raw MercadoLibre item responses with complete product details including variations, attributes, pictures, and pricing.

---

### 2. Transform Phase

**Purpose**: Convert external data formats into Jopit's canonical Item model while preserving data fidelity.

**Implementation Location**: `src/main/domain/utils/meli_transformer.go`

**Main Function**: `TransformMeliItemToJopitItem()` - Orchestrates all transformation steps

**Transformation Steps**:

#### Step 1: Basic Metadata
Extracts fundamental product information like title, description, status, and condition. The description is enriched by prepending the brand name.

**Related Functions**:
- `mapStatus()` - Converts MercadoLibre status values to Jopit's status model
- `mapCategory()` - Maps category IDs and domain IDs to Jopit categories

#### Step 2: Delivery & Packaging
Creates delivery information including package dimensions and fragility settings.

**Related Functions**:
- `mapDelivery()` - Builds the delivery model with dimensions

**Current Limitation**: Uses default dimensions instead of extracting from SELLER_PACKAGE_* attributes.

#### Step 3: Attributes Extraction
Processes product attributes, filtering out invalid or unwanted data.

**Related Functions**:
- `extractAttributes()` - Main attribute processing function
- `ExtractAttributeValue()` - Helper to extract specific attribute values

**What Gets Extracted**:
- Core attributes (gender, composition, color, fabric type)
- Structured MercadoLibre attributes (stored as MeliAttribute objects)
- Sale terms (warranty, returns, payment methods)

**What Gets Filtered Out**:
- Invalid attributes (value_id = "-1" or empty value_name)
- Internal MercadoLibre flags (GIFTABLE, IS_EMERGING_BRAND, etc.)
- Attributes that don't provide value to end users

#### Step 4: Variants Mapping
Groups product variations by color and organizes size/stock information.

**Related Functions**:
- `mapVariants()` - Creates variant array grouped by color
- `extractColorFromVariation()` - Extracts color information
- `extractSizeFromVariation()` - Extracts size labels
- `extractImagesByPictureIDs()` - Maps images to specific variants

**Logic**:
- If no variations exist, creates a single default variant
- Groups variations by color (each color becomes a separate variant)
- For each color variant, extracts all available sizes with stock levels
- Associates images with their corresponding color variant
- Marks the first variant as the main display variant

#### Step 5: Price Extraction
Extracts pricing information and currency details.

**Related Functions**:
- `mapPrice()` - Builds price model
- `mapCurrency()` - Maps currency codes to Jopit's currency objects

**Current Configuration**: Only ARS (Argentine Peso) is active; other currencies are disabled.

#### Step 6: Size Guide Mapping
Converts MercadoLibre size charts into Jopit's size guide format.

**Related Functions**:
- `mapSizeGuide()` - Transforms size chart data
- `ExtractSizeChartID()` - Extracts SIZE_GRID_ID from attributes

**Process**:
1. Extract SIZE_GRID_ID from product attributes
2. Fetch size chart details from MercadoLibre API
3. Parse chart rows and columns into size measurement objects
4. Store external size grid ID for future synchronization

#### Step 7: Source Tracking
Adds comprehensive metadata for ETL provenance and auditing.

**Related Functions**:
- `mapTransformMetadata()` - Creates audit trail metadata

**Tracked Information**:
- Source type (e.g., "meli" for MercadoLibre)
- External ID (original product ID from source system)
- External SKU (seller's SKU)
- Batch ID (groups items from same ETL execution)
- Import timestamp
- ETL version
- Transform metadata (original values, counts, business data)

**Transform Metadata Includes**:
- Original domain and category IDs
- Item counts (variations, pictures, attributes)
- Business data (listing type, buying mode, quantities sold)
- Shipping information
- Timestamps from source system

---

### 3. Load Phase

**Purpose**: Persist transformed items into the Jopit Items domain.

**Implementation Location**: `src/main/domain/services/etl.go` and `src/main/domain/clients/items.go`

**Key Functions**:
- `LoadMercadoLibre()` - Orchestrates full ETL pipeline
- `BulkCreateItems()` - Sends items to Jopit Items API

**Process**:
1. Collect all successfully transformed items
2. Send bulk insert request to Jopit Items API
3. Items API validates, enriches, and stores in MongoDB
4. Track success/failure for each item

**Error Handling**:
- Individual transformation failures don't stop the pipeline
- Failed items are tracked with detailed error information
- Bulk load failures mark all items in that batch as failed
- ETL succeeds if at least one item loads successfully

---

## Error Handling & Resilience

**Implementation Location**: `src/main/domain/services/etl.go`

**Strategies**:

### Individual Item Failures
Each item is transformed independently. If one fails, others continue processing. Failed items are logged with:
- External ID and title
- Stage where failure occurred (transform or load)
- Specific error message

### Panic Recovery
Each transformation is wrapped with defer/recover to catch unexpected panics and convert them to errors.

### Batch Processing
Uses `ETLResult` struct to track:
- Total items processed
- Success count
- Failure count
- Detailed list of failed items

### Success Criteria
The ETL returns success (HTTP 200) if at least one item successfully loads. Only returns error (HTTP 500) if ALL items fail.

---

## Batch Management

**Purpose**: Group items from same ETL execution for tracking and management.

**Batch ID Format**: `{source}-{user_id}` (e.g., "meli-user12345")

**Uses**:
- Track which items came from same ETL run
- Enable batch deletion/rollback
- Audit trail for ETL executions
- Analytics on ETL performance

**Related Endpoints**:
- Delete batch by ID (removes all items from that ETL run)
- Filter items by batch ID

---

## Data Outputs

The ETL process generates three JSON files for debugging and verification:

### 1. `meli-items-extracted.json`
Raw data as received from MercadoLibre API before any transformation.

**Purpose**: Debug extraction issues, verify API responses, understand source data structure.

### 2. `jopit-items-transformed.json`
Full Jopit Item objects after transformation but before loading to database.

**Purpose**: Verify transformation logic, check field mappings, validate data enrichment.

### 3. `jopit-items-failed.json`
List of items that failed during transformation or loading.

**Purpose**: Review failures, identify data quality issues, manual intervention for problematic items.

---

## Field Mappings Summary

### Direct Mappings
- MercadoLibre `title` → Jopit `name`
- MercadoLibre `id` → Jopit `source.external_id`
- MercadoLibre `price` → Jopit `price.amount`
- MercadoLibre `pictures` → Jopit `variants[].images`

### Transformed Mappings
- MercadoLibre `status` (active/paused/closed) → Jopit `status` (active/inactive)
- MercadoLibre `condition` (new/used) → Jopit `attributes.condition` (new/pre-owned)
- MercadoLibre `attributes[BRAND]` + `title` → Jopit `description`

### Structured Mappings
- MercadoLibre `variations[]` → Jopit `variants[]` (grouped by color)
- MercadoLibre `attributes[]` → Jopit `attributes.meli_attributes[]` (with filtering)
- MercadoLibre `sale_terms[]` → Jopit `attributes.meli_attributes[]` (with prefix)

### Metadata Preservation
All transformation metadata is stored in `source.transform_metadata` as key-value pairs for audit and debugging purposes.

---

## Performance Characteristics

### Current Performance
- **Sequential Processing**: Items transformed one at a time
- **Estimated Time**: ~10 seconds for 100 items
- **Size Chart Fetching**: Individual API calls per chart (N+1 pattern)

### Known Bottlenecks
1. No concurrent processing of items
2. Size charts fetched individually rather than in batch
3. All items loaded into memory before processing
4. No pagination for large catalogs (>1000 items)

### Future Optimizations Identified
- Implement worker pool pattern for concurrent transformation
- Batch size chart fetching with caching
- Stream processing for large catalogs
- Add retry logic with exponential backoff

---

## Integration Points

### Authentication
Uses OAuth 2.0 credentials stored in MongoDB (managed by `src/main/domain/repositories/company_layout.go`).

### Items API Client
Communicates with Jopit Items API via REST client (`src/main/domain/clients/items.go`) using bearer token authentication.

### MercadoLibre API
Interacts with MercadoLibre through HTTP client (`src/main/domain/clients/http.go`) with automatic token refresh.

### Database
Credentials stored in MongoDB; transformed items sent to Items API which handles MongoDB persistence.

---

## Monitoring & Observability

### Current Implementation
- OpenTelemetry tracing enabled
- Debug JSON files generated per ETL run
- Detailed error tracking in ETLResult

### Available for Enhancement
- Add Prometheus metrics (items processed, duration, failures)
- Alert on high failure rates
- Dashboard for ETL health monitoring
- Audit log for ETL executions

---

## Related Files Reference

**Core ETL Logic**:
- `src/main/domain/services/etl.go` - Main ETL orchestration
- `src/main/domain/utils/meli_transformer.go` - MercadoLibre transformation logic

**External Communication**:
- `src/main/domain/services/mercadolibre.go` - MercadoLibre service layer
- `src/main/domain/clients/http.go` - MercadoLibre HTTP client
- `src/main/domain/clients/items.go` - Jopit Items API client
- `src/main/domain/clients/shops.go` - Shop validation client

**Data Models**:
- `src/main/domain/models/item.go` - Jopit Item model
- `src/main/domain/models/dto/` - MercadoLibre API response DTOs

**API Endpoints**:
- `src/main/domain/handlers/etl.go` - ETL HTTP handlers

**Configuration**:
- `src/main/api/config/config.go` - Environment configuration

---

**Last Updated**: February 10, 2026  
**ETL Version**: 1.0.0
