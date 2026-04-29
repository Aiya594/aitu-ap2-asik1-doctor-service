# Doctor Service - gRPC Medical Scheduling Platform

## 1. Project Overview

The Doctor Service is the authoritative source for doctor data in a microservices-based Medical Scheduling Platform that uses **gRPC** for all inter-service and client-to-service communication. This service manages the complete lifecycle of doctors in the system and is called by the Appointment Service to validate doctor existence before creating or updating appointments.

The Doctor Service:
- Owns all doctor data and retrieval logic
- Enforces email uniqueness constraint
- Validates all required fields (full_name, email)
- Exposes gRPC RPCs for CRUD operations
- Maintains Clean Architecture principles
- Located in `../aitu-ap2-asik1-doctor-service/` (separate folder from Appointment Service)

UPD:
- **In-memory storage replaced** with a PostgreSQL-backed repository (`internal/repository/repository.go`).
- **Schema managed via migrations** — no DDL in application code; `golang-migrate` runs `migrations/` automatically on startup.
- **NATS publisher added** — after every successful `CreateDoctor`, the service publishes a `doctors.created` event to NATS Core Pub/Sub.
- The publisher is injected behind the `EventPublisher` interface, so broker failures never block the gRPC response.

**Key Role:** Acts as a gRPC dependency for the Appointment Service; called during appointment creation/update to verify doctor existence.

---

---

## Architecture
 
```
gRPC client
     │  CreateDoctor / GetDoctor / ListDoctors
     ▼
transport/grpc  (thin handler, maps errors to gRPC status codes)
     │
     ▼
use-case        (business rules, validation, UUID generation)
     │
     ├──► repository  (PostgreSQL via database/sql + lib/pq)
     │
     └──► event.Publisher  (NATS Core, fire-and-forget)
```
 
---
 

### Service Responsibilities

**Doctor Service:**
- CRUD operations on doctor data via gRPC
- Enforces email uniqueness constraint
- Validates required fields (full_name, email)
- Called by Appointment Service via gRPC client
- No external service dependencies

---

## Start Instruction

Pull  [Notification Service](https://github.com/Aiya594/aitu-ap2-asik3-notification-service) and [Appointment Service](https://github.com/Aiya594/aitu-ap2-asik1-appointment-service) in the same folder with notification folder. Create `docker-compose.yml` according to  `docker-compose.example.yml` and then:

```bash
docker-compose up -d --build
```

---

## Migrations
 
Migration files live in `migrations/`:
 
```
migrations/
├── 000001_create_doctors.up.sql
└── 000001_create_doctors.down.sql
```

## Database Schema
 
Managed exclusively through migration files:
 
```sql
CREATE TABLE doctors (
  id             TEXT        PRIMARY KEY,
  full_name      TEXT        NOT NULL,
  specialization TEXT        NOT NULL DEFAULT '',
  email          TEXT        NOT NULL UNIQUE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
```
 
---
 
## Project Structure
 
```
doctor-service/
├── main.go
├── .env.example
├── Dockerfile
├── go.mod
├── internal/
│   ├── app/           # Wire-up: DB, NATS, repo, use-case, gRPC server
│   ├── config/        # Config + DB connection pool
│   ├── event/         # EventPublisher interface + NATS implementation
│   ├── model/         # Doctor, DoctorCreated event models
│   ├── repository/    # PostgreSQL DoctorRepository
│   ├── transport/
│   │   └── grpc/      # gRPC handler + error mapping
│   └── use-case/      # Business logic
├── migrations/
│   ├── 000001_create_doctors.up.sql
│   └── 000001_create_doctors.down.sql
└── proto/
    ├── doctor.proto
    ├── doctor.pb.go
    └── doctor_grpc.pb.go
```
 
---
 
## Graceful Shutdown
 
The service handles `SIGINT` and `SIGTERM`. On shutdown it:
1. Drains the gRPC server (waits for in-flight RPCs).
2. Closes the NATS connection.
3. Closes the database connection pool.
4. Exits with code 0.