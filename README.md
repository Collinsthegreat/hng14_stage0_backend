# Classify API

## Overview
A production-grade REST API service that predicts the gender of a given name using the Genderize API. Built in Golang using strict Clean Architecture principles.

## Tech Stack
- Go 1.21+
- chi router

## Project Structure
```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── classify.go
│   ├── service/
│   │   └── classify.go
│   ├── model/
│   │   └── classify.go
│   └── client/
│       └── genderize.go
├── pkg/
│   └── response/
│       └── response.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Setup & Run Locally

### Prerequisites
- Go 1.21+

### Clone & Run
```bash
git clone https://github.com/Collinsthegreat/hng14_stage0_backend
cd hng14_stage0_backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

### Environment Variables
| Variable | Default | Description |
|---|---|---|
| PORT | 8080 | Server port |
| GENDERIZE_BASE_URL | https://api.genderize.io | Genderize endpoint |
| HTTP_TIMEOUT_SECONDS | 5 | Outbound HTTP timeout |

## API Reference

### GET /api/classify?name={name}

**Success (200):**
```json
{
  "status": "success",
  "data": {
    "name": "john",
    "gender": "male",
    "probability": 0.99,
    "sample_size": 1234,
    "is_confident": true,
    "processed_at": "2026-04-01T12:00:00Z"
  }
}
```

**Errors:** 400, 422, 500, 502 — all return `{ "status": "error", "message": "..." }`

## Live URL
https://yourapp.domain.app
