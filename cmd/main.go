package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	_ "github.com/joho/godotenv/autoload"

	pb "github.com/lodjim/naboo-email/internal/email"
	"google.golang.org/grpc"
)

var emailFrom string = os.Getenv("EMAIL_ADDRESS")
var emailPassword string = os.Getenv("EMAIL_PWD")
var smtpHost string = os.Getenv("SMTP_HOST")
var smtpPort string = os.Getenv("SMTP_PORT")

type EmailService struct {
	auth      smtp.Auth
	host      string
	port      string
	tlsConfig *tls.Config
	pool      chan *smtp.Client
}

type server struct {
	pb.UnimplementedEmailServer
}

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

func (es *EmailService) getClient() (*smtp.Client, error) {
	select {
	case client := <-es.pool:
		return client, nil
	default:
		conn, err := tls.Dial("tcp", es.host+":"+es.port, es.tlsConfig)
		if err != nil {
			return nil, err
		}
		client, err := smtp.NewClient(conn, es.host)
		if err != nil {
			return nil, err
		}
		if err := client.Auth(es.auth); err != nil {
			return nil, err
		}
		return client, nil
	}
}

func (es *EmailService) releaseClient(client *smtp.Client) {
	es.pool <- client
}

func (es *EmailService) SendEmailToClient(emailFrom, emailTo, subject, body string) error {
	client, err := es.getClient()
	if err != nil {
		return err
	}
	defer es.releaseClient(client)

	if err := client.Mail(emailFrom); err != nil {
		return err
	}

	if err := client.Rcpt(emailTo); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", emailFrom, emailTo, subject, body)

	_, err = wc.Write([]byte(message))
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}
	return nil
}

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
