# autostack-api-inventory

AutoStack Inventory API - Vehicle inventory management and search service.

Part of the AutoStack application family in StackSphere.

## Overview

The `autostack-api-inventory` service provides:
- Vehicle inventory listing and search
- Multi-country vehicle data (US, GB, DE, CA, AU)
- Multi-currency support (USD, GBP, EUR, CAD, AUD)
- RESTful API for vehicle management
- JWT authentication
- Health check endpoints

## Architecture

- **Language:** Go 1.21
- **Port:** 8001
- **Data Storage:** In-memory JSON seed files
- **Authentication:** JWT token-based
- **Database:** None (stateless service using seed data)

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /api/v1/vehicles` - List all vehicles
- `GET /api/v1/vehicles/{id}` - Get vehicle details
- `POST /api/v1/vehicles/search` - Search vehicles with filters

## Local Development

### Prerequisites
- Go 1.21+
- Docker (optional)

### Run Locally

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o server ./cmd/server

# Run
PORT=8001 DATA_PATH=./data/seed ./server
```

### Docker Build

```bash
# Build image
docker build -t autostack-api-inventory:local .

# Run container
docker run -p 8001:8001 autostack-api-inventory:local
```

## Deployment

This component uses the StackSphere repo-per-component model:

- **Workflows:** `.cloudbees/workflows/autostack-api-inventory-ci.yaml`
- **Helm Chart:** `chart/autostack-api-inventory/`
- **Reusable Workflows:** References `stacksphereio/core-workflows`

### Deployment Environments

- **dev:** Non-main branch pushes → `dev` environment
- **integ:** Main branch pushes → `integ` environment

### Database Mode

`no-database` - This service is stateless and uses in-memory seed data.

## Configuration

### Environment Variables

- `PORT` - Service port (default: 8001)
- `DATA_PATH` - Path to seed data files (default: /app/data/seed)
- `JWT_SECRET` - JWT signing secret (required)
- `LOG_LEVEL` - Logging level (default: info)

### Seed Data

Vehicle data is stored in `data/seed/`:
- `vehicles.json` - Vehicle inventory
- `users.json` - Demo user accounts

## Testing

```bash
# Unit tests
go test ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
```

## CloudBees CI/CD

Workflow triggered on:
- All branch pushes
- Manual workflow_dispatch

Pipeline stages:
1. **test** - Run Go tests with Smart Tests integration
2. **build** - Build and push Docker image
3. **deploy-dev** - Deploy to dev environment (non-main branches)
4. **deploy-integ** - Deploy to integ environment (main branch)

## Related Components

- **autostack-api-valuations** - Vehicle valuation service
- **autostack-web** - React UI frontend

## Stack

Part of **AutoStack** - Vehicle marketplace demonstration platform.

**Stack Name:** `autostack`
**Architecture:** StackSphere ARCH v0.5

---

Last updated: 2026-03-10
