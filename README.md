# ECommerce-Platform

A modular, backend e-commerce platform built in Go providing REST and GraphQL APIs, asynchronous event processing, and containerized local infrastructure.

---

## Tech Stack

- **Language:** Go (1.24+)
- **API Frameworks:** Gin (REST), gqlgen (GraphQL)
- **Database & ORM:** PostgreSQL, GORM
- **Database Migrations:** golang-migrate
- **Event Streaming / Queue:** Watermill, AWS SQS (LocalStack)
- **Object Storage:** AWS S3 (LocalStack / AWS SDK v2)
- **Authentication:** JWT (JSON Web Tokens)
- **Documentation:** Swagger / OpenAPI via swaggo
- **Infrastructure:** Docker, Docker Compose, Nginx

---

## Repository Structure

```text
.
├── bin/                 # Compiled binaries
├── cmd/
│   ├── api/             # Main HTTP and GraphQL API entrypoint
│   └── notifier/        # Background worker for processing queue events
├── db/
│   └── migrations/      # Versioned SQL migration files
├── docker/
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── nginx/           # Reverse proxy configuration
├── docs/                # Generated Swagger/OpenAPI documentation
├── graph/               # GraphQL schemas, resolvers, and generated code
├── internal/
│   ├── config/          # Environment configuration loader
│   ├── database/        # DB connection and lifecycle management
│   ├── dto/             # Data transfer objects
│   ├── events/          # Event definitions and publishers
│   ├── models/          # GORM database models
│   ├── notifications/   # Notification handlers
│   ├── provider/        # External service providers (S3, storage)
│   ├── repositories/    # Database query abstractions
│   ├── server/          # HTTP route handlers and middleware
│   └── service/         # Core business logic
├── uploads/             # Local file storage target (when not using S3)
├── Makefile             # Automation commands
├── gqlgen.yml           # gqlgen configuration
└── go.mod
```

---

## Prerequisites

Ensure the following tools are installed on your machine:

- [Go](https://go.dev/dl/) (>= 1.24)
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) (CLI)
- [golangci-lint](https://golangci-lint.run/) (optional, for linting)
- [swag](https://github.com/swaggo/swag) (optional, for regenerating API docs)

---


## Quickstart

### 1. Start Infrastructure Services

Spin up PostgreSQL, LocalStack (S3, SQS), and Nginx using Docker Compose:

```bash
make docker-up
```

### 2. Run Database Migrations

Apply pending database migrations:

```bash
make migrate-up
```

To roll back the migrations:

```bash
make migrate-down
```

To create a new migration pair:

```bash
make migrate-create name=add_table_name
```

### 3. Run the Services

**API Server:**
```bash
make run
```
Or run directly:
```bash
go run ./cmd/api
```

**Notifier Worker:**
In a separate terminal, launch the background consumer:
```bash
go run ./cmd/notifier
```

---

## Build

Compile all service binaries to the `bin/` directory:

```bash
make build
```

This compiles:
- `bin/api`
- `bin/notifier`

---


---

## API Endpoints

Once the API server is running on `http://localhost:8080`:

- **REST API Documentation (Swagger):** `http://localhost:8080/swagger/index.html`
- **GraphQL Playground:** `http://localhost:8080/playground`
- **GraphQL Endpoint:** `http://localhost:8080/query`

---

## License

This project is licensed under the [MIT License](LICENSE).
