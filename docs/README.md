# Jopit ETL Service - Complete Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Features](#features)
4. [Project Structure](#project-structure)
5. [Data Flow](#data-flow)
6. [External Integrations](#external-integrations)
7. [API Endpoints](#api-endpoints)
8. [Configuration](#configuration)
9. [Authentication & Security](#authentication--security)
10. [Data Models](#data-models)
11. [Error Handling](#error-handling)
12. [Deployment](#deployment)
13. [Testing](#testing)
14. [Monitoring & Observability](#monitoring--observability)
15. [Development Guide](#development-guide)
16. [Known Limitations](#known-limitations)
17. [Future Roadmap](#future-roadmap)

---

## Overview

The **Jopit ETL Service** is a specialized microservice responsible for extracting product catalog data from external e-commerce platforms, transforming it into Jopit's standardized data model, and loading it into the Jopit Items domain. It currently supports MercadoLibre as the primary external source, with architecture designed for multi-source extensibility.

### Purpose
- **Data Integration**: Connect Jopit with external e-commerce platforms
- **Catalog Synchronization**: Import seller catalogs from marketplaces into Jopit
- **Data Transformation**: Normalize diverse data formats into Jopit's canonical model
- **Provenance Tracking**: Maintain full audit trail of where data originated

### Key Capabilities
- OAuth 2.0 authentication with automatic token refresh
- Batch processing of product catalogs
- Individual item error handling without pipeline failure
- Size guide mapping from external sources
- Variant and pricing extraction
- Comprehensive metadata preservation
- Rollback capability via batch management

---

## Architecture

### Design Pattern
The service follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────┐
│                   API Layer                         │
│  (Handlers, Routes, Middleware)                     │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│                 Domain Layer                        │
│  • Services (Business Logic)                        │
│  • Repositories (Data Access)                       │
│  • Models (Domain Entities)                         │
│  • Utils (Transformation Logic)                     │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│              Infrastructure Layer                    │
│  • HTTP Clients (External APIs)                     │
│  • Database Clients                                 │
│  • Configuration                                    │
└─────────────────────────────────────────────────────┘
```

### Component Breakdown

**Handlers** (`src/main/domain/handlers/`)
- Receive HTTP requests
- Validate input parameters
- Delegate to services
- Format responses
- Handle HTTP-level errors

**Services** (`src/main/domain/services/`)
- Implement business logic
- Orchestrate ETL pipeline
- Coordinate between multiple clients
- Handle domain-level errors
- Transaction management

**Clients** (`src/main/domain/clients/`)
- Abstract external API communication
- Handle HTTP requests/responses
- Manage authentication
- Connection pooling
- Retry logic

**Repositories** (`src/main/domain/repositories/`)
- Data persistence operations
- MongoDB queries
- Credential management
- BSON serialization

**Utils** (`src/main/domain/utils/`)
- Pure transformation functions
- Data mapping logic
- Validation helpers
- No side effects

**Models** (`src/main/domain/models/`)
- Domain entities (Item, Source, Variant, etc.)
- DTOs for external APIs
- Validation rules
- JSON/BSON tags

---

## Features

### Current Features (Version 1.0.0)

#### ✅ MercadoLibre Integration
- Full OAuth 2.0 authentication flow
- Automatic token refresh (proactive 1-hour before expiry)
- Batch item fetching
- Size chart retrieval
- Comprehensive attribute extraction

#### ✅ ETL Pipeline
- Extract: Fetch raw data from external sources
- Transform: Convert to Jopit's Item model
- Load: Bulk insert to Jopit Items API
- Per-item error handling with graceful degradation

#### ✅ Data Transformation
- Status and condition mapping
- Category classification
- Variant grouping by color
- Size/stock extraction
- Price and currency conversion
- Attribute filtering (removes invalid data)
- Sale terms extraction
- Size guide mapping

#### ✅ Provenance Tracking
- Source type identification
- External ID preservation
- Batch ID for grouping
- Import timestamps
- ETL version tracking
- Complete transform metadata

#### ✅ Error Management
- Individual item failure tracking
- Panic recovery
- Detailed error messages
- Failed items JSON export
- Partial success handling

#### ✅ Batch Operations
- Batch deletion by batch ID
- Group operations on ETL runs
- Rollback capability

#### ✅ Observability
- OpenTelemetry distributed tracing
- Debug JSON file outputs
- Structured logging
- HTTP request/response tracking

---

## Project Structure

```
jopit-api-etl/
├── docs/
│   ├── README.md                    # This file
│   ├── etl-process.md              # Detailed ETL process explanation
│   └── etl-project-review.md       # Quality assessment (7.2/10)
│
├── environment/
│   ├── docker-compose.yml          # Local development setup
│   └── jopit-api-items.dockerfile  # Container definition
│
├── src/
│   ├── main/
│   │   ├── api/
│   │   │   ├── main.go            # Application entry point
│   │   │   ├── app/
│   │   │   │   ├── app.go         # Gin server setup
│   │   │   │   └── routes.go      # Route definitions
│   │   │   ├── config/
│   │   │   │   └── config.go      # Environment configuration
│   │   │   ├── dependencies/
│   │   │   │   ├── builder.go     # Dependency injection setup
│   │   │   │   └── manager.go     # Dependency lifecycle
│   │   │   └── platform/
│   │   │       └── storage/       # Database connections
│   │   │
│   │   └── domain/
│   │       ├── clients/
│   │       │   ├── http.go        # MercadoLibre HTTP client
│   │       │   ├── items.go       # Jopit Items API client
│   │       │   └── shops.go       # Jopit Shops API client
│   │       │
│   │       ├── handlers/
│   │       │   ├── company_layout.go  # CSV ETL endpoints
│   │       │   └── etl.go            # MercadoLibre ETL endpoints
│   │       │
│   │       ├── models/
│   │       │   ├── company_layout.go  # CSV schema models
│   │       │   ├── item.go           # Jopit Item model
│   │       │   ├── prices.go         # Price model
│   │       │   ├── shop.go           # Shop model
│   │       │   └── dto/              # External API DTOs
│   │       │
│   │       ├── repositories/
│   │       │   └── company_layout.go  # Credential repository
│   │       │
│   │       ├── services/
│   │       │   ├── company_layout.go  # CSV ETL service
│   │       │   └── etl.go            # MercadoLibre ETL service
│   │       │
│   │       └── utils/
│   │           ├── utils.go          # General utilities
│   │           └── meli_transformer.go  # MercadoLibre transformation
│   │
│   └── tests/
│       └── internal/              # Unit and integration tests
│
├── go.mod                         # Go module definition
├── go.sum                         # Dependency checksums
├── README.md                      # Project README
├── etl.go                         # Utility script
└── meli-response.json            # Sample API response
```

---

## Data Flow

### Complete MercadoLibre ETL Flow

```
User Request
    ↓
[ETL Handler] (etl.go)
    ↓
[ETL Service] (etl.go)
    ↓
┌─────────────── EXTRACT ───────────────┐
│                                        │
│  [MercadoLibre Service]               │
│     ↓                                  │
│  [Credentials Repository]             │
│     ↓                                  │
│  [HTTP Client] (with auto-refresh)    │
│     ↓                                  │
│  MercadoLibre API                     │
│     ↓                                  │
│  Raw Item Data []                     │
│                                        │
└────────────────────────────────────────┘
    ↓
┌─────────────── TRANSFORM ─────────────┐
│                                        │
│  For each item:                       │
│     ↓                                  │
│  [Meli Transformer] (meli_transformer.go) │
│     ├── Extract attributes            │
│     ├── Map variants (color/size)     │
│     ├── Map price & currency          │
│     ├── Extract size guide            │
│     ├── Build source metadata         │
│     └── Create Jopit Item             │
│                                        │
│  Error handling per item              │
│                                        │
└────────────────────────────────────────┘
    ↓
┌─────────────── LOAD ──────────────────┐
│                                        │
│  [Items Client] (items.go)            │
│     ↓                                  │
│  Jopit Items API                      │
│     ↓                                  │
│  MongoDB (via Items API)              │
│                                        │
└────────────────────────────────────────┘
    ↓
[ETL Result]
    ├── Success count
    ├── Failure count
    └── Failed items details
```

---

## External Integrations

### 1. MercadoLibre API

**Base URL**: `https://api.mercadolibre.com`

**Endpoints Used**:
- `GET /users/{user_id}/items/search` - List seller's items
- `GET /items?ids=ID1,ID2,ID3` - Batch fetch item details
- `GET /size_charts/{chart_id}` - Get size guide data
- `POST /oauth/token` - Token refresh

**Authentication**: OAuth 2.0 with refresh tokens

**Rate Limits**: Not explicitly handled (future enhancement)

**Documentation**: https://developers.mercadolibre.com

### 2. Jopit Items API

**Base URL**: Configured via `ITEMS_API_URL` environment variable

**Endpoints Used**:
- `POST /items/bulk` - Bulk create items
- `DELETE /items/batch/{batch_id}` - Delete items by batch

**Authentication**: Bearer token (Firebase JWT)

**Features Used**: Schema validation, data enrichment, MongoDB persistence

### 3. Jopit Shops API

**Base URL**: Configured via environment

**Endpoints Used**:
- `GET /shops/{shop_id}` - Validate shop existence

**Authentication**: Bearer token (Firebase JWT)

### 4. Firebase Authentication

**Purpose**: JWT token validation for API security

**Integration**: Middleware validates bearer tokens on all protected routes

**User Context**: User ID extracted from JWT for multi-tenant operations

### 5. MongoDB

**Purpose**: Credential storage

**Database**: Configured via `DB_NAME` environment variable

**Collections**:
- `company_layout` - Stores MercadoLibre OAuth credentials
- Indexed by user_id and shop_id

**Driver**: Official MongoDB Go driver

---

## API Endpoints

### MercadoLibre ETL Endpoints

#### `POST /etl/mercadolibre/load`
Load products from MercadoLibre into Jopit.

**Headers**:
- `Authorization: Bearer {firebase_jwt}`

**Request Body**:
```json
{
  "user_id": "firebase_user_id",
  "shop_id": "jopit_shop_id"
}
```

**Response** (200 OK):
```json
{
  "batch_id": "meli-firebase_user_id",
  "total_items": 50,
  "success_count": 48,
  "failure_count": 2,
  "failed_items": [
    {
      "external_id": "MLA123456",
      "title": "Product that failed",
      "failure_stage": "transform",
      "error_message": "invalid price format"
    }
  ]
}
```

**Response** (500 Internal Server Error):
```json
{
  "message": "all items failed to load",
  "status": "etl_failed",
  "code": 500,
  "causes": ["detailed error messages"]
}
```

#### `DELETE /etl/batch/{batch_id}`
Delete all items from a specific ETL batch.

**Headers**:
- `Authorization: Bearer {firebase_jwt}`

**Path Parameters**:
- `batch_id` - The batch identifier (e.g., "meli-user123")

**Response** (200 OK):
```json
{
  "message": "batch deleted successfully",
  "batch_id": "meli-user123"
}
```

### MercadoLibre OAuth Endpoints

#### `GET /mercadolibre/auth-url`
Get MercadoLibre authorization URL for OAuth flow.

**Query Parameters**:
- `user_id` - Firebase user ID
- `shop_id` - Jopit shop ID

**Response** (200 OK):
```json
{
  "auth_url": "https://auth.mercadolibre.com/authorization?response_type=code&client_id=..."
}
```

#### `POST /mercadolibre/callback`
Handle OAuth callback and exchange code for tokens.

**Request Body**:
```json
{
  "code": "TG-...",
  "user_id": "firebase_user_id",
  "shop_id": "jopit_shop_id"
}
```

**Response** (200 OK):
```json
{
  "message": "credentials saved successfully"
}
```

### CSV ETL Endpoints

#### `POST /csv/load`
Load items from CSV file upload.

**Headers**:
- `Authorization: Bearer {firebase_jwt}`
- `Content-Type: multipart/form-data`

**Form Data**:
- `file` - CSV file
- `user_id` - User ID
- `shop_id` - Shop ID
- `layout_id` - Schema configuration ID

---

## Configuration

### Environment Variables

Required configuration (set in `.env` or deployment environment):

```bash
# Server
PORT=8080
GIN_MODE=release  # or "debug" for development

# Database
DB_USERNAME=mongouser
DB_PASSWORD=mongopass
DB_CLUSTER=cluster.mongodb.net
DB_NAME=jopit-etl

# MercadoLibre OAuth
MERCADO_LIBRE_CLIENT_ID=1234567890
MERCADO_LIBRE_CLIENT_SECRET=your_secret_here
MERCADO_LIBRE_REDIRECT_URI=https://your-domain.com/callback

# External APIs
ITEMS_API_URL=https://items-api.jopit.com
SHOPS_API_URL=https://shops-api.jopit.com

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=https://tempo.grafana.com
OTEL_EXPORTER_OTLP_HEADERS=Authorization=Basic xyz...
```

### Configuration Loading

Configuration is loaded via `src/main/api/config/config.go` using environment variables with fallbacks.

---

## Authentication & Security

### OAuth 2.0 Flow (MercadoLibre)

**Step 1: Authorization URL**
- User clicks "Connect MercadoLibre"
- Frontend requests auth URL from ETL API
- ETL API builds authorization URL with client ID and redirect URI
- Frontend redirects user to MercadoLibre

**Step 2: User Authorization**
- User logs into MercadoLibre
- User approves permissions
- MercadoLibre redirects back with authorization code

**Step 3: Token Exchange**
- Frontend sends authorization code to ETL API callback
- ETL API exchanges code for access token and refresh token
- Tokens stored in MongoDB (encrypted at rest)

**Step 4: Token Usage**
- ETL API reads tokens from MongoDB before API calls
- Includes access token in Authorization header
- Automatically refreshes if expiring within 1 hour

### JWT Authentication (Jopit APIs)

All ETL endpoints protected by Firebase JWT middleware:
- Validates bearer token signature
- Extracts user ID from token claims
- Ensures user can only access their own data
- Returns 401 Unauthorized for invalid tokens

### Security Best Practices Implemented

- ✅ No tokens in frontend/localStorage
- ✅ OAuth tokens stored in backend database
- ✅ Automatic token refresh (no manual intervention)
- ✅ JWT validation on all protected routes
- ✅ MongoDB credentials encrypted at rest
- ✅ HTTPS required in production
- ✅ User-scoped data access (multi-tenant isolation)

---

## Data Models

### Core Item Model

The Jopit Item model represents a product in the catalog. Located in `src/main/domain/models/item.go`.

**Key Structures**:
- **Item**: Top-level product entity
- **Source**: ETL provenance tracking
- **Attributes**: Product attributes including MercadoLibre-specific data
- **Variant**: Color variant with size/stock information
- **SizeGuide**: Measurement information
- **Price**: Pricing and currency details
- **Delivery**: Shipping and packaging information

All fields use `omitempty` JSON tags to avoid sending null/empty values.

### Source Tracking Model

Preserves complete ETL metadata:
- `source_type`: Origin platform (e.g., "meli")
- `external_id`: Original product ID from source
- `external_sku`: Seller's SKU
- `batch_id`: Groups items from same ETL run
- `imported_at`: ETL execution timestamp
- `synced_at`: Last incremental sync (future use)
- `etl_version`: ETL code version
- `transform_metadata`: Key-value audit data

### MercadoLibre-Specific Models

**MeliAttribute** structure for preserving raw attribute data:
- `id`: Attribute ID (e.g., "BRAND")
- `name`: Display name (e.g., "Marca")
- `value_id`: Value identifier
- `value_name`: Display value (e.g., "Nike")

Stored in `attributes.meli_attributes` array.

---

## Error Handling

### Error Hierarchy

**API Errors** (`apierrors` package):
- Structured error responses
- HTTP status codes
- Error causes array
- Consistent format across all endpoints

**Domain Errors**:
- Business logic validation errors
- Data transformation errors
- External API errors

**Infrastructure Errors**:
- Database connection errors
- HTTP client errors
- Configuration errors

### ETL-Specific Error Handling

**Individual Item Failures**:
- Transformation errors caught per item
- Failed items logged with details
- Processing continues for remaining items

**Panic Recovery**:
- Each transformation wrapped with defer/recover
- Panics converted to regular errors
- Stack trace logged for debugging

**Partial Success**:
- ETL succeeds if ≥1 item loads successfully
- Detailed breakdown in response (success/failure counts)
- Failed items array with specific error messages

**Complete Failure**:
- Returns 500 error only if ALL items fail
- Includes aggregated error causes
- Developer-friendly error messages

---

## Deployment

### Docker Support

Dockerfile located at `environment/jopit-api-items.dockerfile` (note: filename needs update).

**Build**:
```bash
docker build -f environment/jopit-api-items.dockerfile -t jopit-etl:latest .
```

**Run**:
```bash
docker run -p 8080:8080 --env-file .env jopit-etl:latest
```

### Docker Compose

Local development setup at `environment/docker-compose.yml`.

**Start services**:
```bash
docker-compose -f environment/docker-compose.yml up
```

### Cloud Deployment

Recommended for production:
- **Container Orchestration**: Kubernetes or Cloud Run
- **Secrets Management**: Google Secret Manager or Vault
- **Load Balancing**: Cloud Load Balancer
- **Auto-scaling**: Based on CPU/memory metrics
- **Health Checks**: `/health` endpoint

---

## Testing

### Current Test Coverage

**Status**: Limited test coverage (~4/10 rating from project review)

**Existing Tests**:
- Located in `src/tests/internal/`
- Setup utilities for mock data

**Critical Gap**: Core ETL transformation logic (`meli_transformer.go`) has zero tests.

### Running Tests

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Recommended Test Priorities

1. **Unit Tests** for `meli_transformer.go` (highest priority)
2. **Integration Tests** for full ETL pipeline
3. **Contract Tests** for MercadoLibre API responses
4. **Error Handling Tests** for failure scenarios

---

## Monitoring & Observability

### Current Implementation

**OpenTelemetry Tracing**:
- Distributed tracing enabled
- Traces sent to Grafana Cloud
- Request/response spans
- External API call tracing

**Debug Outputs**:
- `meli-items-extracted.json` - Raw source data
- `jopit-items-transformed.json` - Transformed items
- `jopit-items-failed.json` - Failure details

**Structured Logging**:
- Gin request logging
- Error logging with context
- ETL execution logs

### Recommended Enhancements

**Metrics** (not yet implemented):
- `etl_items_processed_total` - Counter
- `etl_duration_seconds` - Histogram
- `etl_failures_total` - Counter by stage
- `etl_batch_size` - Gauge

**Alerts**:
- Failure rate >10%
- ETL duration >5 minutes
- Token refresh failures
- API rate limit hits

**Dashboards**:
- ETL execution timeline
- Success/failure breakdown
- Performance metrics
- Error trends

---

## Development Guide

### Prerequisites

- Go 1.23.9 or higher
- MongoDB instance
- MercadoLibre developer account
- Firebase project

### Local Setup

1. Clone repository:
```bash
git clone <repository-url>
cd jopit-api-etl
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment:
```bash
cp .env.example .env
# Edit .env with your credentials
```

4. Run locally:
```bash
go run src/main/api/main.go
```

### Adding a New ETL Source

To add support for another platform (e.g., Shopify):

1. **Create DTO models** in `src/main/domain/models/dto/`
2. **Implement HTTP client** in `src/main/domain/clients/`
3. **Create transformation utility** in `src/main/domain/utils/`
4. **Add service methods** in `src/main/domain/services/`
5. **Create handlers** in `src/main/domain/handlers/`
6. **Update routes** in `src/main/api/app/routes.go`
7. **Add tests** in `src/tests/`

### Coding Standards

- Follow Clean Architecture principles
- Keep handlers thin (delegate to services)
- Make utils pure functions (no side effects)
- Use dependency injection
- Write tests for all transformation logic
- Document public functions
- Use meaningful variable names
- Handle errors explicitly (no silent failures)

---

## Known Limitations

### Current Constraints

1. **Performance** (6.5/10):
   - Sequential processing (no concurrency)
   - ~10 seconds for 100 items
   - N+1 queries for size charts

2. **Testing** (4/10):
   - Core transformer untested
   - No integration tests for full pipeline

3. **Observability** (5/10):
   - No metrics/alerts
   - No dashboards
   - Limited structured logging

4. **Functionality Gaps**:
   - No pagination (fails at >1000 items)
   - No incremental sync
   - No retry logic
   - No rate limiting

5. **Data Quality**:
   - Hardcoded packaging dimensions
   - Color hex codes default to black
   - Category mapping uses heuristics
   - No variant-specific pricing

6. **Security**:
   - No API rate limiting
   - No audit log
   - No data encryption in transit validation

---

## Future Roadmap

### High Priority (Next Sprint)

1. **Implement Unit Tests**
   - Target: 80% coverage for transformation logic
   - Priority: `meli_transformer.go`

2. **Add Pagination**
   - Handle large catalogs (>1000 items)
   - Configurable page size

3. **Concurrent Processing**
   - Worker pool pattern (10 workers)
   - Expected: 10x speedup

4. **Retry Logic**
   - Exponential backoff
   - Max 3 retries for transient failures

### Medium Priority (Next Quarter)

5. **Prometheus Metrics**
   - Items processed, duration, failures
   - Integration with Grafana

6. **Incremental Sync**
   - Delta sync using timestamps
   - Only process changed items

7. **Shared Models Module**
   - Extract Item model to separate package
   - Share between ETL and Items API

8. **Enhanced Transform Metadata**
   - Extract packaging dimensions
   - Add 15+ business fields

### Low Priority (Future)

9. **Additional ETL Sources**
   - Shopify integration
   - WooCommerce integration
   - CSV improvements

10. **Rate Limiting**
    - Respect external API limits
    - Implement backoff strategies

11. **Advanced Features**
    - Event-driven architecture
    - Real-time sync
    - Conflict resolution

---

## Support & Contributing

### Documentation

- **ETL Process Details**: See [etl-process.md](etl-process.md)
- **OAuth 2.0 Flow**: See [oauth-flow.md](oauth-flow.md)
- **Project Review**: See [etl-project-review.md](etl-project-review.md)
- **API Documentation**: Swagger/OpenAPI (to be implemented)

### Getting Help

For issues or questions:
- Check existing documentation
- Review debug JSON outputs
- Check application logs
- Contact development team

### Contributing

When contributing:
1. Follow coding standards
2. Write tests for new features
3. Update documentation
4. Submit pull request with description
5. Ensure CI/CD passes

---

## Version History

**Version 1.0.0** (Current)
- MercadoLibre OAuth integration
- Complete ETL pipeline
- Source tracking
- Error handling
- Batch operations
- OpenTelemetry tracing

---

**Project Status**: Production-ready for small-medium catalogs (<500 items)  
**Quality Score**: 7.2/10 (see [project review](etl-project-review.md) for details)  
**Last Updated**: February 10, 2026
