# System Architecture

## Overview

This SaaS API manages real estate properties, tenants, and associated operations.
It is written in Go, follows RESTful API principles, and is containerized using Docker.

## Key Components

### 1. API Server (`cmd/server`)

- Entry point for HTTP server
- Initializes routes, middleware, and services

### 2. Application Layer (`internal/`)

- Business logic and service orchestration
- Follows clean architecture principles
- Separated into domains (e.g., `property`, `tenant`, `user`)

### 3. Persistence Layer (`internal/repository`)

- Interfaces for data access
- Uses PostgreSQL via standard `database/sql` or an ORM (if applicable)

### 4. API Routes (`api/`)

- HTTP endpoints organized by resource
- Uses standard net/http or third-party router (e.g., chi, gorilla/mux)

### 5. Configuration (`configs/`)

- Environment variable definitions
- Configuration structs parsed on startup

## Data Flow

Client → HTTP Request → Router → Handler → Service Layer → Repository → Database

## Technologies

| Component | Stack                  |
| --------- | ---------------------- |
| Language  | Go                     |
| API       | REST (JSON)            |
| DB        | PostgreSQL             |
| CI/CD     | GitHub Actions         |
| Runtime   | Docker                 |
| Hosting   | DigitalOcean (Planned) |

## Future Enhancements

- GraphQL API
- Background jobs (e.g., notifications)
- Audit logging
