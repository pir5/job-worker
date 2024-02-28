# README for Health Worker

## Overview
Health Worker is a comprehensive solution designed to monitor and manage the health status of services by performing various health checks and applying routing policies based on the health check outcomes. It integrates seamlessly with Echo framework for middleware and utilizes Echo's context, GORM for database interactions, and Redis for managing state and configurations. This service is capable of handling health checks over HTTP(S) and TCP, enabling dynamic routing policy adjustments for high availability and fault tolerance.

## Getting Started

### Prerequisites
- Go 1.20+
- MySQL or compatible databases for storing health check configurations and routing policies
- Redis for managing failed health check counters and worker configurations
- Echo framework for handling HTTP requests and middleware

### Installation
Clone the repository and install the necessary dependencies.

```sh
git clone https://github.com/path/to/health-worker.git
cd health-worker
go install ./...
```

### Configuration
Edit the configuration file (`/etc/health-worker/worker.toml` by default) to set up the database, Redis, and other operational parameters like the server port, polling intervals, and concurrency settings.

### Running Health Worker
Health Worker consists of three main components: Register Server, Worker Server, and API Server. Each component is started with its respective command.

#### Start Register Server
To start the register server which handles the registration of health checks and scheduling:

```sh
health-worker register --config /path/to/your/config.toml
```

#### Start Worker Server
To start the worker server which executes the health checks:

```sh
health-worker worker --config /path/to/your/config.toml
```

#### Start API Server
To start the API server which provides an interface for managing health checks and routing policies:

```sh
health-worker api --config /path/to/your/config.toml
```

## Features

### Health Checks
Define health checks with various parameters like interval, threshold, and the type of check (HTTP, HTTPS, TCP) along with specific parameters for each type.

### Routing Policies
Define routing policies that determine how traffic should be routed based on the outcome of health checks. Supports simple, failover, and detach types.

### Dynamic Configuration
Utilize environment variables and a configuration file to dynamically adjust worker, database, and Redis settings.

### Middleware Integration
Use Echo's middleware capabilities to easily integrate header-based authentication for API requests.

### Extensible Design
Easily extend the health check and routing policy models for custom logic and integrations.

## Documentation
Swagger-based API documentation is automatically generated and accessible through the API server, providing a detailed reference for all available endpoints and their usage.

## Contributing
Contributions are welcome! Please feel free to submit pull requests or open issues to discuss new features or improvements.

## License
Health Worker is open-source software licensed under the MIT license.
