# Devices Management REST API (Golang Version)

A high-performance microservice designed in Go 1.23+ for managing device lifecycles and operational states. This project decouples business rule restrictions from core data access layers using idiomatic HTTP Middlewares, matching the enterprise standards previously established in the Java Spring iteration.

---

## 🛠️ Tech Stack & Key Choices

- **Language:** Go 1.23+ (Statical compilation, fast startup, lightweight footprint).
- **Web Framework:** [Gin Gonic](https://github.com) (High-performance routing triage engine).
- **Database Driver:** [pgx/v5](https://github.com) (Native PostgreSQL connection pool handling with zero reflection overhead).
- **Database Engine:** PostgreSQL 15 (Relational persistent storage).
- **Documentation:** [swag](https://github.com) (Automatic OpenAPI 3.0 specification generator via code comments).
- **Testing Toolkit:** [testify/suite](https://github.com) & [testcontainers-go](https://github.com) (Advanced database integration testing).

---

## 📐 Project Architecture Layout

The codebase strictly adheres to a clean, decoupled layer separation mapping architecture:

```text
devices-api-go/
├── config/       # Database connection pool allocation and initialization
├── controller/   # REST API endpoint multiplexing handlers and test scenarios
├── docs/         # Automatically compiled Swagger OpenAPI static definitions
├── middleware/   # Centralized HTTP interceptors (Centralized error handler & business rules)
├── model/        # Entities, custom types (pseudo-enums), and DTO schemas
└── repository/   # Raw SQL statements with efficient offset pagination handling
```

---

## 🔒 Implemented Business Safeguards

The HTTP Middleware layer dynamically intercepts modification inputs before they ever propagate down to the persistence service tier, enforcing strict evaluation text rules:
1. **Creation Time Constraint:** Returns `400 Bad Request` with the text `"Creation time cannot be updated."` if the client attempts to alter immutable metadata timestamps via PATCH operations.
2. **Resource Elimination Safeguard:** Returns `400 Bad Request` with the text `"In use devices cannot be deleted."` if an HTTP DELETE target is currently associated with an `IN_USE` operational state.
3. **Property Modification Lock:** Returns `400 Bad Request` with the text `"Name and brand properties cannot be updated if the device is in use."` if a client tries to modify the `name` or `brand` attributes of an actively running asset.

---

## 🚀 Getting Started (Docker Compose Execution)

To spin up both the isolated PostgreSQL database instance and the Go microservice instantly without local engine dependencies, execute the following command at the root directory:

```bash
docker compose up --build -d
```

Once the containers shift into a running state:
- The REST API will accept live traffic on port `8080`.
- The PostgreSQL cluster will map safely to host port `5432`.

To gracefully terminate the execution environment cluster, run:
```bash
docker compose down
```

---

## 📖 Live Documentation (Swagger UI)

When the application environment is online via Docker, you can inspect and interact with the full live interactive API endpoints by pointing your web browser to:

👉 **`http://localhost:8080/swagger/index.html`**

---

## 🧪 Running Automated Tests

The testing suite simulates full context lifecycles to guarantee the validation shields are sound:
- **Controller Pack:** Validates HTTP router state flows, validation tags, and global error middleware handling.
- **Repository Pack (Advanced Integration):** Utilizes **`testcontainers-go`** to dynamically spin up a true, transient sandboxed PostgreSQL container instance inside Docker at runtime to validate live infrastructure database compliance metrics.

To execute all integration scenarios and inspect outcomes across all packages, trigger the native testing engine via your terminal console:

```bash
go test -v ./...
```

---

## 🔮 Future Production-Readiness Improvements

Per the structural roadmap of this microservice, the following items outline the architectural enhancements and future scalability improvements:

1. **True Total Counts Pagination Metadata**: The current baseline `GetAll` endpoint outputs a high-performance raw JSON array slice to maximize read speeds on millions of rows. In a full production layout, a lightweight parallel `SELECT COUNT(*)` query or an estimated row count counter would be bundled via a wrapper struct to return total pages to the client application.
2. **Database Migrations Engine**: Database schemas are currently handled natively. Integrating a tool like `golang-migrate` or `Liquibase` would allow structural version tracking directly inside the repository pipeline.
3. **Structured Context Logging & Metrics**: Replacing the standard `log` package with a structured JSON logger like `uber-go/zap` or Go's native `slog` would improve log parsing. Exposing a `/metrics` endpoint via Prometheus would enable real-time dashboard observability.
