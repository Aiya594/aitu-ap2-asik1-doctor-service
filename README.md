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

**Key Role:** Acts as a gRPC dependency for the Appointment Service; called during appointment creation/update to verify doctor existence.

---

## 2. Scope and Constraints

### What Changes

- **All HTTP/REST endpoints are replaced with gRPC endpoints**
  - Gin HTTP server replaced with gRPC server
  - `transport/http/` replaced with `transport/grpc/`

- **Service exposes a gRPC server**
  - Listens on port 50051 for gRPC connections
  - No REST functionality

- **Protocol Buffers define the service contract**
  - `.proto` file committed to the repository
  - Generated Go stubs (`*pb.go`, `*_grpc.pb.go`) committed alongside


---

## 3. Architecture Overview

### Service Responsibilities

**Doctor Service:**
- CRUD operations on doctor data via gRPC
- Enforces email uniqueness constraint
- Validates required fields (full_name, email)
- Called by Appointment Service via gRPC client
- No external service dependencies

### Project Structure

```
aitu-ap2-asik1-doctor-service/
├── main.go                         # Entry point
├── go.mod / go.sum
├── proto/                          # Protocol Buffer definitions
│   ├── doctor.proto                # Doctor service contract
│   ├── doctor.pb.go                # Generated stubs (committed)
│   └── doctor_grpc.pb.go           # Generated gRPC stubs (committed)
├── internal/
│   ├── model/
│   │   └── doctor.go               # Domain model (unchanged from A1)
│   ├── repository/
│   │   ├── repository.go           # Storage interface
│   │   └── errors.go               # Repository errors
│   ├── use-case/
│   │   ├── doctror_use_case.go     # Business logic (unchanged from A1)
│   │   └── errors.go               # Use-case errors
│   ├── transport/
│   │   ├── grpc/
│   │   │   └── server.go           # gRPC server handlers (REPLACES http/)
│   │   └── http/                   # REMOVED (replaced by grpc/)
│   └── app/
│       └── app.go                  # Application setup (gRPC server instead of Gin)
└── README.md                       # This file
```

### Dependency Flow (Clean Architecture)

```
gRPC Handler
    ↓
UseCase (Business Logic)
    ↓
Repository (In-Memory Storage)
```

**Rules:**
- gRPC handler unmarshal proto messages, calls use case, returns proto responses
- Use case contains all business logic; imports NO protobuf types
- Repository handles storage only
- Mapping between proto messages and domain models happens ONLY in the gRPC layer



---

## 6. Installing and Regenerating Proto Stubs

### Prerequisites

1. **Install protoc** (Protocol Buffer Compiler)
   - Windows: Download from [protobuf releases](https://github.com/protocolbuffers/protobuf/releases) and add to PATH

2. **Install Go gRPC plugins**
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

### Regenerate Stubs

From the `aitu-ap2-asik1-doctor-service/` directory:

```bash
protoc --go_out=. --go-grpc_out=. proto/doctor.proto
```


---

## 7. Running the Doctor Service


### Startup

```bash
cd aitu-ap2-asik1-doctor-service
go run main.go
```

**gRPC Port:** `localhost:50051`

The service should be started **before** the Appointment Service because Appointment Service will connect to it during startup or on first request.
