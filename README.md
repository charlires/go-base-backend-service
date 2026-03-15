# go-base-backend-service

## Project Overview

This project is a Go-based backend service that implements a simple music playlist management system. It allows users to create playlists, add tracks to playlists, and retrieve playlist details. The service is structured using clean architecture principles, with clear separation of concerns between the core domain logic and the input/output adapters.

```plaintext
├── main.go
├── internal
│   ├── adapters
│   │   ├── input
│   │   │   └── playlist_http_adapter.go
│   │   └── output
│   │       └── playlist_sqlite_repository.go
│   │       └── ...
│   └── core
│       ├── domain
│       │   ├── playlist_domain.go
│       │   └── ...
│       └── services
│       │   ├── playlist_service.go
│       │   └── ...
│       └── ports
│           └── repository.go
├── init.sql
├── go.mod
└── go.sum
```

## Project Setup

Populate the database using the provided SQL script:
```bash
sqlite3 database.db -init init.sql 
```

Install dependencies:
```bash
go mod tidy
go mod vendor
```

Run the application:
```bash
go run main.go
```