# Music Analyze

A small Go service that exposes analytics endpoints on the Chinook music dataset. It uses MySQL for storage, chi for routing, and Swagger (swaggo) for interactive API docs.

## Features
- Top genres by sales with percentage of total
- Top countries by customers and by average order value
- Sales by media type with percentage of total
- Highest-sale invoices with customer names
- Invoice statistics (median and mode)
- Sales totals by support agents

## Tech Stack
- Go 1.22
- Router: `github.com/go-chi/chi/v5`
- MySQL driver: `github.com/go-sql-driver/mysql`
- Swagger UI: `github.com/swaggo/http-swagger` + generated docs in `./docs`

## Quick Start (Local)
1. Configure MySQL connection in `config.yaml` (or copy `config.docker.yaml` and adjust):
   - `server.port` default is `5000`
   - Set `database.host`, `port`, `user`, `password`, `name`
2. Run the API:
   - `go run ./cmd`
   - Or specify a custom config: `CONFIG_PATH=config.docker.yaml go run ./cmd`
3. Visit Swagger UI: `http://localhost:5000/swagger/index.html`

## Swagger Docs (generate/update)
This repo already includes generated docs in `./docs`. If you change annotations and need to regenerate:
- Install the CLI: `go install github.com/swaggo/swag/cmd/swag@latest`
- Generate docs: `$(go env GOPATH)/bin/swag init -g cmd/main.go -o docs`

You can also run a Swagger-only mode without a DB connection to preview the docs:
- `CONFIG_PATH=config.swagger.yaml SWAGGER_ONLY=1 go run ./cmd`
- Open `http://localhost:5051/swagger/index.html`

## Endpoints
- `GET /analytics/top-genres?limit=10`
- `GET /analytics/top-countries?limit=5`
- `GET /analytics/media-types`
- `GET /analytics/highest-sale`
- `GET /analytics/invoice-stats`
- `GET /analytics/agents`

See Swagger for detailed schemas and examples.

## Docker
Build and run the containerized API:
- Build: `docker build -t music-analyze:latest .`
- Run: `docker run --rm -p 5000:5000 -e CONFIG_PATH=/app/config.yaml music-analyze:latest`

Compose (simple app container):
- `docker-compose up --build`
- API listens on `http://localhost:5000`

Note: You must provide a reachable MySQL database via configuration. The sample `config.docker.yaml` expects a host named `db`, but this compose file does not start MySQL. Point to your own DB or extend the compose file with a MySQL service.

## Database
- A Chinook schema/data script is provided at `internal/infra/db/mysql/DB.sql`.
- The queries in this service are written for a lowercase/snake_case variant of Chinook. If you import the original CamelCase schema, you may need to adjust the table/column names or use the provided script to create a compatible schema.

## Configuration
- `CONFIG_PATH`: path to the YAML config (default `config.yaml`)
- `SWAGGER_ONLY`: if `1` or `true`, starts only the Swagger UI without DB connections (useful for previewing docs)

## Project Layout
- `cmd/main.go`: entrypoint, router, Swagger wiring
- `internal/http/analytics`: HTTP handlers and routes
- `internal/application/analyst`: application service layer
- `internal/domain/analyst`: domain entities and repository port
- `internal/infra/db/mysql`: MySQL repository implementation and connection
- `docs/`: generated Swagger files (imported by the app)

## Development
- Tidy deps: `go mod tidy`
- Build: `go build ./...`
- Run: `go run ./cmd`

## License
This project includes the Chinook sample database script. See the original Chinook license for terms. The rest of the code is provided without a specific license unless otherwise stated.
