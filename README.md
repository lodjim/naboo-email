
# Naboo Email Server

The Naboo Email Server is a gRPC-based Go application for securely sending emails via SMTP with TLS encryption. It incorporates an SMTP client connection pooling mechanism to enhance performance, reliability, and scalability. Ideal for public use, the project is open-source on GitHub.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Setup & Configuration](#setup--configuration)
- [Running the Server](#running-the-server)
- [API Usage](#api-usage)
- [Code Structure](#code-structure)
- [Future Enhancements](#future-enhancements)
- [Contributing](#contributing)
- [License](#license)

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

1. **Clone the Repository:**

```bash
git clone https://github.com/yourusername/naboo-email.git
cd naboo-email
```

2. **Install Dependencies:**

```bash
go mod tidy
```

3. **Build and Run the Server:**

```bash
cd cmd
go build -o naboo-email-server
./naboo-email-server
```

The server will start listening on `0.0.0.0:50051` and log its status.

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

## Future Enhancements

- Improved logging and monitoring.
- Enhanced error handling to prevent abrupt termination of the server.
- Expanded connection pooling management for high-throughput scenarios.

## Contributing

We welcome contributions! Fork the repo, implement your features or fixes, and submit pull requests. Please discuss significant changes via issues before submitting PRs.

## License

This project is licensed under the [MIT License](LICENSE).
