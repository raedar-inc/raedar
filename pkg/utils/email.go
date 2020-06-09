package utils

import (
	"net/smtp"
	"os"
)

type Email struct {
	Email string
}

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// address URI to smtp server
func (s *smtpServer) address() string {
	return s.host + ":" + s.port
}

func (e Email) SendEmail(emailSubject, msgBody string) error {
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	senderEmail := os.Getenv("SENDER_EMAIL_ADDRESS")
	senderPass := os.Getenv("SENDER_EMAIL_PASSWORD")
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpServer.host)

	to := []string{e.Email}
	message := []byte("To: " + e.Email + "\r\n" +
		"Subject: \r\n" + emailSubject + "\r\n" +
		"\r\n" +
		msgBody + "\r\n")

	err := smtp.SendMail(smtpServer.address(), auth, senderEmail, to, message)
	if err != nil {
		return err
	}
	return nil
}
