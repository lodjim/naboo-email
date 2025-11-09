
# Naboo Email Server

[![Docker Image](https://img.shields.io/docker/v/cifope/naboo-email?label=docker&logo=docker)](https://hub.docker.com/r/cifope/naboo-email)
[![Docker Image Size](https://img.shields.io/docker/image-size/cifope/naboo-email/latest)](https://hub.docker.com/r/cifope/naboo-email)
[![Docker Pulls](https://img.shields.io/docker/pulls/cifope/naboo-email)](https://hub.docker.com/r/cifope/naboo-email)
[![License](https://img.shields.io/github/license/lodjim/naboo-email)](LICENSE)

The Naboo Email Server is a gRPC-based Go application for securely sending emails via SMTP with TLS encryption. It incorporates an SMTP client connection pooling mechanism to enhance performance, reliability, and scalability. Ideal for public use, the project is open-source on GitHub.

**Key Features:**
- 🚀 Ultra-lightweight Docker image (~15MB)
- 🔒 TLS 1.2+ encryption for all SMTP connections
- ⚡ Connection pooling for optimal performance
- 🌍 Multi-platform support (linux/amd64, linux/arm64)
- 📦 Ready-to-use Docker images on Docker Hub

## Table of Contents

- [Quick Start with Docker](#quick-start-with-docker)
- [Overview](#overview)
- [Architecture](#architecture)
- [Setup & Configuration](#setup--configuration)
- [Running the Server](#running-the-server)
- [API Usage](#api-usage)
- [Code Structure](#code-structure)
- [CI/CD & Versioning](#cicd--versioning)
- [Future Enhancements](#future-enhancements)
- [Contributing](#contributing)
- [License](#license)

## Quick Start with Docker

The fastest way to get started is using our pre-built Docker images:

```bash
# Pull the latest image
docker pull cifope/naboo-email:latest

# Run the server
docker run -d \
  --name naboo-email \
  -p 50051:50051 \
  -e EMAIL_ADDRESS="your-email@example.com" \
  -e EMAIL_PWD="your-password" \
  -e SMTP_HOST="smtp.example.com" \
  -e SMTP_PORT="465" \
  -e POOL_SIZE="10" \
  cifope/naboo-email:latest
```

Available tags:
- `latest` - Latest stable release
- `1.0.0` - Specific version
- `1.0` - Latest patch of minor version 1.0
- `1` - Latest minor/patch of major version 1

## Overview

Naboo Email Server is built in Go and leverages gRPC for efficient remote procedure calls, exposing methods to send emails securely. Configuration is straightforward using environment variables.

## Architecture

The application comprises:

- **gRPC Server**: Implements the `SendEmail` method for sending emails via API requests.
- **EmailService**: Handles SMTP client management, including:
  - TLS-secured connections to SMTP servers
  - SMTP authentication
  - Connection pooling for efficient resource management
- **SMTP Integration**: Uses Go's built-in `net/smtp` and `crypto/tls` packages.

## Setup & Configuration

Create a `.env` file in the project's root directory with the following required variables:

```bash
EMAIL_ADDRESS="your-email@example.com"
EMAIL_PWD="your-email-password"
SMTP_HOST="smtp.example.com"
SMTP_PORT="465"
POOL_SIZE="10"  # Optional, defaults to 5 if not set
```

## Running the Server

### Option 1: Using Docker (Recommended)

See [Quick Start with Docker](#quick-start-with-docker) above.

### Option 2: Building from Source

1. **Clone the Repository:**

```bash
git clone https://github.com/lodjim/naboo-email.git
cd naboo-email
```

2. **Install Dependencies:**

```bash
go mod tidy
```

3. **Build and Run the Server:**

```bash
cd cmd
go build -o naboo-email-server main.go
./naboo-email-server
```

The server will start listening on `0.0.0.0:50051` and log its status.

### Option 3: Building Docker Image Locally

```bash
# Build the image
docker build -t naboo-email:local .

# Run the container
docker run -d \
  --name naboo-email \
  -p 50051:50051 \
  --env-file .env \
  naboo-email:local
```

## API Usage

### Proto Definition

```proto
message SendEmailRequest {
    string email_target = 1;
    string subject = 2;
    string message = 3;
}

message SendEmailReply {
    string message = 1;
}
```

### Client Example

```go
client.SendEmail(ctx, &pb.SendEmailRequest{
    EmailTarget: "recipient@example.com",
    Subject:     "Greetings from Naboo",
    Message:     "<h1>Hello!</h1><p>This is a test email from Naboo Email Server.</p>",
})
```

### Regenerate gRPC Code

Run this to regenerate the Go code from `.proto` files:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/email/emailservice.proto
```

## Code Structure

- **Main (`main.go`):** Initializes environment variables, creates the SMTP client pool, and starts the gRPC server.
- **EmailService:**
  - Manages SMTP client connections with TLS encryption.
  - Includes functions for acquiring and releasing SMTP clients from a connection pool.
  - Handles SMTP commands (`MAIL FROM`, `RCPT TO`, `DATA`) to send emails.

## CI/CD & Versioning

This project uses GitHub Actions for automated Docker image building and publishing to Docker Hub.

### Automated Builds

- **On push to `main`**: Builds and pushes images with commit SHA tags
- **On version tags (`v*`)**: Builds and pushes semantic version tags (e.g., v1.0.0 → 1.0.0, 1.0, 1, latest)
- **Security scanning**: All images are scanned with Trivy for vulnerabilities

### Image Optimization

- Multi-stage build with scratch base image
- Image size: ~15MB (98% reduction from standard Go images)
- Multi-platform support: linux/amd64, linux/arm64
- Non-root user execution for enhanced security
- Automated build caching for faster CI/CD

### Creating a New Release

```bash
# Create and push a new version tag
git tag -a v1.1.0 -m "Release version 1.1.0"
git push origin v1.1.0

# GitHub Actions will automatically:
# - Build multi-platform Docker images
# - Run security scans
# - Push to Docker Hub with appropriate tags
# - Update Docker Hub description
```

## Future Enhancements

- Improved logging and monitoring
- Enhanced error handling to prevent abrupt termination of the server
- Expanded connection pooling management for high-throughput scenarios
- Metrics and observability (Prometheus, OpenTelemetry)

## Contributing

We welcome contributions. Fork the repo, implement your features or fixes, and submit pull requests. Please discuss significant changes via issues before submitting PRs.

## License

This project is licensed under the [MIT License](LICENSE).
