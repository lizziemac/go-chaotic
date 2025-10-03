# Basic Go HTTP Server

This repository provides a minimal Go HTTP server with built-in logging and a simple ping/pong endpoint. It includes basic tests and can be used as a template for quickly starting new Go server projects without re-creating the directory structure/setup every time.

## Layout

```
├── cmd/
│   └── api/           # Entrypoint
│       └── main.go
├── internal/
│   └── api/           # API-specific middleware and utils
│       └── ...
│   └── middleware/    # Shared middleware
│       └── ...
│   └── logger/        # Shared logger
│       └── ...
├── go.mod
├── go.sum
└── README.md
```

## Running the Server

Start the server:

```sh
go run cmd/api/main.go
```

## Example Requests

### Happy Path

Send a GET request to the ping endpoint:

```sh
curl localhost:8080/api/v1/ping
```

**Expected Response:**

```json
{"message":"pong"}
```

**Expected Server Logs:**

```
time=2025-10-02T16:13:27.730-04:00 level=INFO msg="request received" method=GET path=/api/v1/ping body=""
time=2025-10-02T16:13:27.730-04:00 level=INFO msg="request complete" method=GET path=/api/v1/ping duration=165.75µs
```

### Unhappy Path

Request an invalid endpoint:

```sh
curl localhost:8080/bad
```

**Expected Response:**

```json
{"error":"not found"}
```

**Expected Server Log:**

```
[WARN] GET /bad 404
```

## Running Tests

It is recommended to run tests in your IDE if possible, but to run manually:

```sh
go test -v ./...
```

To run with a clean cache:

```sh
go test -count=1 -v ./...
```

## Generating the Documentation

To browse Go documentation for the `internal/api` package locally, run:

```sh
go doc -http
```

This command starts a local documentation server and automatically opens it. Once opened, you can navigate to the package to view detailed documentation for all exported symbols.

## Design Considerations

### State Management in Middleware

Originally, the project was going to use context injection for state management, but that risked hiding important details in Context and encouraging scope creep.

Instead, it now uses a custom HTTP handler that explicitly passes state between middleware. This introduces some boilerplate and doesn't prevent SharedState from becoming a God struct, but it feels more transparent and less "sneaky".