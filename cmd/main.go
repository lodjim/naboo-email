package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/lodjim/naboo-email/internal/email"
)

var (
	emailFrom     = mustGetEnv("EMAIL_ADDRESS")
	emailPassword = mustGetEnv("EMAIL_PWD")
	smtpHost      = mustGetEnv("SMTP_HOST")
	smtpPort      = mustGetEnv("SMTP_PORT")
)

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

type EmailService struct {
	auth      smtp.Auth
	host      string
	port      string
	tlsConfig *tls.Config
	pool      chan *smtp.Client
	timeout   time.Duration
}

type server struct {
	pb.UnimplementedEmailServer
	emailService *EmailService
}

func NewEmailService(host, port, emailFrom, emailPassword string, poolSize int) *EmailService {
	auth := smtp.PlainAuth("", emailFrom, emailPassword, host)
	tlsConfig := &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	}

	pool := make(chan *smtp.Client, poolSize)
	return &EmailService{
		auth:      auth,
		host:      host,
		port:      port,
		tlsConfig: tlsConfig,
		pool:      pool,
		timeout:   10 * time.Second,
	}
}

func (es *EmailService) getClient() (*smtp.Client, error) {
	select {
	case client := <-es.pool:
		if err := client.Noop(); err == nil {
			return client, nil
		}
		client.Close()
		return es.createNewClient()
	default:
		return es.createNewClient()
	}
}

func (es *EmailService) createNewClient() (*smtp.Client, error) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: es.timeout},
		"tcp",
		net.JoinHostPort(es.host, es.port),
		es.tlsConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("TLS dial failed: %w", err)
	}

	client, err := smtp.NewClient(conn, es.host)
	if err != nil {
		return nil, fmt.Errorf("SMTP client creation failed: %w", err)
	}

	if err := client.Auth(es.auth); err != nil {
		client.Close()
		return nil, fmt.Errorf("SMTP auth failed: %w", err)
	}

	return client, nil
}

func (es *EmailService) releaseClient(client *smtp.Client) {
	if err := client.Reset(); err != nil {
		client.Close()
		return
	}

	select {
	case es.pool <- client:
	default:
		client.Close()
	}
}

func (es *EmailService) SendEmailToClient(emailFrom, emailTo, subject, body string) error {
	client, err := es.getClient()
	if err != nil {
		return fmt.Errorf("failed to get SMTP client: %w", err)
	}
	defer es.releaseClient(client)

	// Set timeout for SMTP operations
	ctx, cancel := context.WithTimeout(context.Background(), es.timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		if err := client.Mail(emailFrom); err != nil {
			done <- fmt.Errorf("MAIL FROM failed: %w", err)
			return
		}

		if err := client.Rcpt(emailTo); err != nil {
			done <- fmt.Errorf("RCPT TO failed: %w", err)
			return
		}

		wc, err := client.Data()
		if err != nil {
			done <- fmt.Errorf("DATA command failed: %w", err)
			return
		}

		message := fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nSubject: %s\r\n"+
				"MIME-version: 1.0\r\n"+
				"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
			emailFrom, emailTo, subject, body,
		)

		if _, err := wc.Write([]byte(message)); err != nil {
			done <- fmt.Errorf("message write failed: %w", err)
			return
		}

		if err := wc.Close(); err != nil {
			done <- fmt.Errorf("message close failed: %w", err)
			return
		}

		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("SMTP operation timed out")
	}
}

func (s *server) SendEmail(ctx context.Context, in *pb.SendEmailRequest) (*pb.SendEmailReply, error) {
	err := s.emailService.SendEmailToClient(emailFrom, in.EmailTarget, in.Subject, in.Message)
	if err != nil {
		log.Printf("Email send failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to send email: %v", err)
	}

	log.Printf("Email sent to %s", in.EmailTarget)
	return &pb.SendEmailReply{Message: "Email sent successfully"}, nil
}

func main() {
	poolSize, err := strconv.Atoi(os.Getenv("POOL_SIZE"))
	if err != nil {
		poolSize = 5
		log.Printf("Using default pool size: %d", poolSize)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	emailService := NewEmailService(smtpHost, smtpPort, emailFrom, emailPassword, poolSize)
	s := grpc.NewServer()
	pb.RegisterEmailServer(s, &server{emailService: emailService})

	log.Printf("Server running on :50051 with %d SMTP connections in pool", poolSize)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
