# doctor-service


## Project Overview

The Doctor Service manages doctor data. It acts as a source of truth for doctors in the system.

## Purpose

This service may be used to encapsulate doctor-related data and support other services ([Appointment Service](https://github.com/Aiya594/aitu-ap2-asik1-appointment-service) )

It is structured separately to enforce data ownership, enable independent scaling, avoid cross-service coupling.

## Service Responsibilities

The Doctor Service owns:
- Doctor data storage
- Doctor retrieval
- Doctor existence validation

Typical operations:
- Add a new doctor
- Get doctor by ID
- List doctors

## Folder Structure 

(example — similar structure expected)
```
internal/
├── transport/http/ -> HTTP handlers
├── use-case/       -> Business logic
├── repository/     -> Storage layer
├── model/          -> Domain models
├── app/            -> Server management
```

## Dependency Flow

```
Handler -> UseCase -> Repository
```

No external dependencies.

# How to Run
1. Set environment variables as given in ```.env.example```

2. Run service
```
go run main.go
```

3. Test endpoint
- ```POST /doctors``` - add a doctor too storage
```
Example:
{
    "full_name":"Test Doctor",
    "email":"doctor@example.com",
    "specialization":"surgery",
}
```

- ```GET /doctors/{id}``` - return a doctor by ID

- ```GET /doctors``` - returns all doctors in storage


## Why No Shared Database?

The Doctor Service owns all doctor data.

Other services must not access its database and only use its API instead

This ensures strict boundaries, independent evolution and safe data management
