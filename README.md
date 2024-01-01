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