# Structured Chaos Reverse Proxy Server

This repository provides a minimal reverse proxy built in Go for injecting chaos (latency, dropped connections) into API requests. It acts as an intermediary layer to test system resilience, allowing you to dynamically configure chaos rules per user via HTTP headers.

**By default, this proxy does not apply any rules, and is simply a passthrough. You must configure per-user rules to
see the bad behavior**.

Inspired by khizar-sudo's [chaos-proxy](https://github.com/khizar-sudo/chaos-proxy)

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

⚠️ Before starting the proxy server, you need a backend server running for it to forward requests to.

Set the `TARGET_URL` environment variable to point to your backend (e.g., a local Go server):

```sh
export TARGET_URL="http://127.0.0.1:8080"
go run cmd/api/main.go
```

## How to Use

### 💨 Passthrough

Send a GET request to the ping endpoint of your backend, but using the address of the proxy server.

```sh
curl localhost:8081/api/v1/ping
```

**Expected Response:**
You should get a normal response back, as if you sent the request directly to your backend.

```json
{"message":"pong"}
```

### 🚶🏻‍♀️Passthrough with User

Send a GET request to the ping endpoint of your backend, but using the address of the proxy server.

```sh
curl -H 'X-User-ID: lizziemac' localhost:8081/api/v1/ping
```

**Expected Response:**
You should get a normal response back, as if you sent the request directly to your backend, IF no configurations are set for your user. If you have latency mode set, all responses will be delayed by the specified duration. If you have drop mode set, responses will be dropped with the specified rate. See [Setting Proxy Configs](#setting-proxy-configs) for more details.

```json
{"message":"pong"}
```

### Setting Proxy Configs

* `mode` (integer): A bitmapped value that determines what chaos to apply. (Note: Because this is a bitmap, a mode of 3 applies both Drop and Delay rules).
    * 0 (PassMsg): Do nothing; pass the message normally.
    * 1 (DropMsg): Drop incoming messages based on the drop_rate.
    * 2 (DelayMsg): Delay messages by the latency_delay.

* `drop_rate` (float, optional): The percentage rate at which messages are dropped (used if mode includes DropMsg). For example, 0.9 drops 90% of the requests.

* `latency_delay_ns` (integer, optional): The amount of time to delay a message in nanoseconds (used if mode includes DelayMsg). For example, 5000000000 equals 5 seconds.

You can dynamically configure rules for specific users by providing an `X-User-ID` header. For example, to inject a 5-second latency delay with a 50% drop rate (mode: 3 represents DelayMsg & DropMsg), and the delay is in nanoseconds:

```sh
curl -X PUT http://localhost:8081/proxy/api/v1/config \
    -H "X-User-ID: test-user" \
    -H "Content-Type: application/json" \
    -d '{"mode": 3, "drop_rate": 0.5, "latency_delay_ns": 5000000000}'
```

**Expected Response:**
```json
{
    "mode": 3,
    "drop_rate": 0.5,
    "latency_delay_ns": 5000000000,
    "ttl": "2026-06-07T16:35:06.22257-04:00"
}
```

### Getting Proxy Configs
To retrieve the current chaos configuration for a specific user, use a GET request:

```sh
curl -X GET http://localhost:8081/proxy/api/v1/config \
    -H "X-User-ID: test-user" 
```

**Expected Response:**
```json
{
    "mode": 1,
    "drop_rate": 0.5,
    "latency_delay_ns": null,
    "ttl": "2026-06-07T16:39:38.139464-04:00"
}
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

### User-Based Configurations using Global Memory

This project uses dependency injection to pass around a global store of configurations. This store is intentionally generic so that the in-memory solution can be replaced with a Redis cache or something else as desired with relative ease.

### Making this a Template Project
This project in its current state can be helpful for sandboxing. However, if you wanted to use it formally in a dev environment, you may have a different way you want to separate configs (e.g. a different header/unique identifier), add specific behaviors, etc. Or you may want to use Redis or something instead of the server's heap. Either way, this was a fun learning project for me, and I hope you benefit from using it!
