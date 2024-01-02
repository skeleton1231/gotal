# README for GoTAL Project

## Introduction
Welcome to GoTAL, a Go-based enterprise project. This project utilizes a robust architecture suitable for scalable and efficient web services.

## Getting Started

### Prerequisites
- Golang environment setup
- MySQL and Redis services running
- Necessary environment variables set for MySQL and Redis configurations

### Project Structure
```
gotal/
├── api
├── build
├── cert
├── cmd
│   └── apiserver
├── configs
├── deployments
├── docs
├── githooks
├── go.mod
├── go.sum
├── init
├── internal
│   ├── apiserver
│   └── pkg
├── logs
├── pkg
├── test
└── tools
```

### Setup and Running
1. **Building the API server:**
   Navigate to the API server directory:
   ```
   cd /Users/huanghaitao/gotal/cmd/apiserver
   ```
   Build the API server:
   ```
   go build
   ```

2. **Starting the Server:**
   Run the API server with the specified configuration:
   ```
   ./apiserver --config ../../configs/apiserver.yaml
   ```
   This will initialize the server with various configurations as shown in your provided start-up log.
   
## Full Configuration for the GoTAL API Server

After the section on starting the server, the following detailed configurations are applied:

### RESTful Service Configuration
```yaml
server:
  mode: debug # Modes: release, debug, test. Default is release.
  healthz: true # Enable health check, setting up /healthz route. Default is true.
  middlewares: recovery,logger,secure,nocache,cors,dump # List of gin middlewares.
```

### gRPC Service Configuration
```yaml
grpc:
  bind-address: 0.0.0.0 # IP address for gRPC. Default is 0.0.0.0.
  bind-port: 8082 # Port for gRPC. Default is 8081.
```

### HTTP Configuration (Insecure)
```yaml
insecure:
  bind-address: 0.0.0.0 # IP address for insecure binding. Default is 127.0.0.1.
  bind-port: 8884 # Non-secure port. Default is 8080.
```

### HTTPS Configuration (Secure)
```yaml
secure:
  bind-address: 0.0.0.0 # IP address for HTTPS. Default is 0.0.0.0.
  bind-port: 8445 # Port for HTTPS. Default is 8443.
  tls:
    cert-key:
      cert-file: /Users/huanghaitao/gotal/cert/apiserver.pem # Certificate file path.
      private-key-file: /Users/huanghaitao/gotal/cert/apiserver-key.pem # Private key file path.
```

### MySQL Database Configuration
```yaml
mysql:
  host: 127.0..0.1 # MySQL server address.
  username: root # MySQL username.
  password:  # MySQL password.
  database: db # Database name.
  max-idle-connections: 100 # Max idle connections.
  max-open-connections: 100 # Max open connections.
  max-connection-life-time: 10s # Connection lifetime.
  log-level: 4 # Log level.
```

### Redis Configuration
```yaml
redis:
  host: 127.0.0.1 # Redis host.
  port: 6379 # Redis port.
  password:  # Redis password.
  # Additional configuration details can be specified here.
```

### JWT Configuration
```yaml
jwt:
  realm: JWT # JWT realm identifier.
  key:  # Secret key.
  timeout: 24h # Token expiration time.
  max-refresh: 24h # Token refresh time.
```

### Feature Configuration
```yaml
feature:
  enable-metrics: true # Enable metrics at /metrics.
  profiling: true # Enable performance analysis at /debug/pprof/.
```

### Rate Limiting Configuration
```yaml
ratelimit:
  requests-per-second: 1 # Requests per second per user.
  burst-size: 20 # Maximum burst size.
  custom-limits:
    "/test-response":
      requests-per-second: 1.5 # Requests per second for specific endpoint.
      burst-size: 10 # Burst size for specific endpoint.
```

### Logging Configuration
```yaml
log:
  name: apiserver # Logger name.
  development: true # Development mode.
  level: debug # Log level.
  format: console # Log format.
  enable-color: true # Color output.
  disable-caller: false # Caller information.
  disable-stacktrace: false # Stack trace.
  output-paths: /Users/huanghaitao/gotal/logs/apiserver.log # Output paths.
  error-output-paths: /Users/huanghaitao/gotal/logs/apiserver.error.log # Error log paths.
```

This detailed configuration will ensure that your GoTAL API server is set up with the specific settings required for its operation. These settings include server modes, service bindings, database connections, logging, and more, ensuring a comprehensive and robust setup for your enterprise-grade application.

### Configuration
The server's behavior is controlled by various flags and options. These include but are not limited to:
- **gRPC Configurations**: Address, port, message size limits.
- **MySQL and Redis Configurations**: Host, port, authentication, pooling settings.
- **JWT Settings**: Key, timeout, refresh settings.
- **Server Mode**: Including debug, test, and release modes.

Detailed flag descriptions are available in the start-up log section of this README.

## Architecture Overview

The project is structured into several layers to promote separation of concerns and maintainability:

1. **Controller Layer**: Handles HTTP requests, invoking the appropriate services.
2. **Service Layer**: Contains business logic and interacts with the repository layer.
3. **Repository Layer**: Responsible for data access and storage management.

### Key Components
- **Cobra & Viper**: Used for building CLI and managing configuration files.
- **Validator**: Ensures that incoming requests meet the defined constraints.
- **Middleware**: Includes authentication, logging, CORS, rate limiting, etc.
- **Logging**: Structured logging for tracing and monitoring.
- **MySQL & Redis**: Used for data storage and caching.

## Contribution
Refer to `CONTRIBUTING.md` for guidelines on how to contribute to this project.

## License
This project is licensed under the terms mentioned in `LICENSE`.

## Changelog
For a detailed changelog, see `CHANGELOG`.

---

This README provides a basic overview of GoTAL. For detailed documentation, refer to the `docs` directory.