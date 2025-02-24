# Naboo Email Server

The Naboo Email Server is a gRPC-based service that sends emails via an SMTP server. It leverages TLS for secure connections and includes a simple SMTP client pooling mechanism to improve performance. This project is designed for public use via GitHub.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Setup & Configuration](#setup--configuration)
- [Code Walkthrough](#code-walkthrough)
  - [Main Entry Point](#main-entry-point)
  - [EmailService](#emailservice)
  - [gRPC Service](#grpc-service)
- [Usage](#usage)
- [Run Server](#run-server)
- [Error Handling & Future Enhancements](#error-handling--future-enhancements)
- [Contributing](#contributing)
- [License](#license)

## Overview

The Naboo Email Server is built in Go and exposes a gRPC API for sending emails. It uses environment variables for configuration, allowing you to easily integrate your SMTP settings. The server listens on TCP port `50051` and accepts gRPC requests defined in the protocol buffer files.

## Architecture

- **gRPC Server**: Exposes the `SendEmail` method, accepting a request that contains the email target, subject, and message.
- **EmailService**: Manages SMTP connections including:
  - Creating TLS-secured connections to the SMTP server.
  - Authenticating using the provided email credentials.
  - Maintaining a pool of SMTP clients to reuse existing connections.
- **SMTP Integration**: Utilizes Go's `net/smtp` and `crypto/tls` packages to send emails securely.

## Setup & Configuration

Before running the server, you must create a `.env` file in `cmd` folder and add these variables:

- `EMAIL_ADDRESS`: The email address to send from.
- `EMAIL_PWD`: The password for the sender's email account.
- `SMTP_HOST`: The hostname of your SMTP server.
- `SMTP_PORT`: The port for the SMTP server.

### Example

```bash
EMAIL_ADDRESS="your-email@example.com"
EMAIL_PWD="your-email-password"
SMTP_HOST="smtp.example.com"
SMTP_PORT="465"
```

### Build & Run

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/naboo-email.git
   cd naboo-email
   ```

2. **Download Dependencies**

   ```bash
   go mod tidy
   ```

3. **Build the Project**

   ```bash
   go build -o naboo-email-server
   ```

4. **Run the Server**

   ```bash
   ./naboo-email-server
   ```

   The server will start listening on `0.0.0.0:50051`.

## Code Walkthrough

### Main Entry Point

The `main` function is the starting point of the server:

- It sets up a TCP listener on port `50051`.
- It creates and registers a gRPC server.
- It binds the email service implementation to the gRPC server.

```go
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterEmailServer(s, &server{})

	fmt.Println("Server is runnnig on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

### EmailService

The `EmailService` struct is responsible for connecting to the SMTP server, authenticating, and sending emails. It also implements a simple connection pooling mechanism.

#### Key Components:

- **Fields**:
  - `auth`: SMTP authentication using `smtp.PlainAuth`.
  - `host` & `port`: SMTP server details.
  - `tlsConfig`: TLS configuration to ensure secure connections.
  - `pool`: A channel that acts as a pool for `smtp.Client` instances.

- **Functions**:
  - `NewEmailService`: Initializes an `EmailService` with the SMTP settings and a pool size.
  - `getClient`: Retrieves an SMTP client from the pool or creates a new one if none is available.
  - `releaseClient`: Returns a used client back to the pool.
  - `SendEmailToClient`: Constructs and sends the email using SMTP commands.

Example snippet from the `EmailService`:

```go
func NewEmailService(host, port, emailFrom, emailPassword string, poolSize int) *EmailService {
	auth := smtp.PlainAuth("", emailFrom, emailPassword, host)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}
	pool := make(chan *smtp.Client, poolSize)
	return &EmailService{
		auth:      auth,
		host:      host,
		port:      port,
		tlsConfig: tlsConfig,
		pool:      pool,
	}
}
```

### gRPC Service

The gRPC service is implemented in the `server` struct, which includes the `SendEmail` method. This method:

1. Receives a `SendEmailRequest` containing the target email, subject, and message.
2. Instantiates an `EmailService` using the environment variables.
3. Calls `SendEmailToClient` to send the email.
4. Returns a `SendEmailReply` with a status message.

```go
func (s *server) SendEmail(ctx context.Context, in *pb.SendEmailRequest) (*pb.SendEmailReply, error) {
	emailTo := in.EmailTarget
	subject := in.Subject
	body := in.Message
	var emailService *EmailService = NewEmailService(smtpHost, smtpPort, emailFrom, emailPassword, 5)
	err := emailService.SendEmailToClient(emailFrom, emailTo, subject, body)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	fmt.Println("Email sent successfully")

	return &pb.SendEmailReply{Message: "sent"}, nil
}
```

#### Usage

Once the server is running, clients can send a gRPC request to the `SendEmail` method. The request is expected to conform to the protocol buffer definition similar to:

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

#### Regenerate pb files

To regenerate **`pb.go`** needed to communicate via gRPC, run the following command :

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       internal/email/emailservice.proto
```

## Run Server

To run the server, move into `cmd` folder and launch the program.
```bash
cd cmd
go run .
```

A sample client might send the following request:

```go
client.SendEmail(ctx, &pb.SendEmailRequest{
    EmailTarget: "recipient@example.com",
    Subject:     "Hello from Naboo Email Server",
    Message:     "<h1>This is a test email</h1>",
})
```

## Error Handling & Future Enhancements

- **Error Handling**:
  The server currently logs a fatal error (`log.Fatalf`) if sending an email fails, which terminates the process. Future improvements might include more graceful error handling to prevent service downtime.

- **Future Enhancements**:
  - Replace `log.Fatalf` with non-terminating error handling.
  - Add more granular logging and monitoring.
  - Expand connection pooling to handle high email throughput.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](LICENSE).
