# Devices Management REST API (Golang Version)

A high-performance microservice designed in Go 1.23+ for managing device lifecycles and operational states. This project decouples business rule restrictions from core data access layers using an idiomatic Aspect-Oriented Programming (AOP) Decorator pattern at the service layer, matching the enterprise standards previously established in the Java Spring iteration.

---

## 🛠️ Tech Stack & Key Choices

- **Language:** Go 1.23+ (Statical compilation, fast startup, lightweight footprint).
- **Pagination:** High-performance Cursor-Based Pagination designed for high-scale enterprise environments.
- **Web Framework:** [Gin Gonic](https://github.com) (High-performance routing triage engine).
- **Database Driver:** [pgx/v5](https://github.com) (Native PostgreSQL connection pool handling with zero reflection overhead).
- **Database Engine:** PostgreSQL 16 (Relational persistent storage).
- **Documentation:** [swag](https://github.com) (Automatic OpenAPI 3.0 specification generator via code comments).
- **Testing Toolkit:** [testify/suite](https://github.com) (Advanced behavioral assertion testing).

---

## 📐 Project Architecture Layout

The codebase strictly adheres to a clean, decoupled layer separation mapping architecture matching Clean Architecture standards:

```text
devices-api-go/
├── cmd/
│   └── api/      # Application entrypoint (main.go with AOP weaving)
├── config/       # Database connection pool allocation and initialization
├── docs/         # Automatically compiled Swagger OpenAPI static definitions
└── internal/     # Private application layers
    ├── domain/   # Entities, custom types, DTO schemas, and Interfaces
    ├── handler/  # REST API HTTP controller handlers (DeviceController)
    ├── middleware/ # Global error mapping interceptors
    ├── repository/ # Database persistence layer (DeviceRepository implementation)
    └── service/  # Core business services and the AOP aspect decorator proxy
```

---

## 🔒 Implemented Business Safeguards (AOP Layer)

The core business rule validation is fully encapsulated inside a generic service-level **AOP Aspect Proxy Decorator** (`DeviceServiceAspect`). It dynamically intercepts modification inputs before they ever propagate down to the persistence service tier, enforcing strict evaluation text rules:
1. **Creation Time Constraint:** Returns `400 Bad Request` with the text `"Creation time cannot be updated."` if the client attempts to alter immutable metadata timestamps via PUT/PATCH operations.
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
- The PostgreSQL cluster will map safely to local host port `5433` (preventing local conflicts with default port 5432).

To gracefully terminate the execution environment cluster, run:
```bash
docker compose down --volumes
```

---

## 📖 Live Documentation (Swagger UI)

When the application environment is online via Docker, you can inspect and interact with the full live interactive API endpoints by pointing your web browser to:

👉 **`http://localhost:8080/swagger/index.html`**

---

## 🧪 Running Automated Tests

The testing suite simulates full context lifecycles to guarantee the validation shields are sound:
- **Service & AOP Aspect Pack:** Validates core domain logic execution boundaries and validation rules using decoupled mock interface repositories.
- **Handler/HTTP Pack:** Simulates network requests and captures JSON payloads, testing global error mapping behaviors.

To execute all integration scenarios and inspect outcomes across all packages, trigger the native testing engine via your terminal console:

```bash
go test -v ./...
```

---

## 🔮 Future Production-Readiness Improvements

Per the structural roadmap of this microservice, the following items outline the architectural enhancements and future scalability improvements:

1. **Database Migrations Engine**: Database schemas are currently handled natively via container initialization scripts. Integrating a tool like `golang-migrate` would allow structural version tracking directly inside the repository pipeline.
2. **Structured Context Logging & Metrics**: Replacing the standard `log` package with a structured JSON logger like Go's native `slog` would improve log parsing. Exposing a `/metrics` endpoint via Prometheus would enable real-time dashboard observability.
