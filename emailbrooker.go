package itswizard_m_emailserviceBrooker

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

type Email struct {
	SmtpServer string
	Port       uint
	Password   string
	Username   string
}

func (p *Email) SendEmail(fromEmailAdress string, toEmailAdress string, EmailSubject string, EmailContent string) (err error) {
	var log string
	from := mail.Address{"", fromEmailAdress}
	to := mail.Address{"", toEmailAdress}
	subj := "This is the email subject"
	body := "This is an example body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := fmt.Sprint(p.SmtpServer, ":", p.Port)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", p.Username, p.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log = log + fmt.Sprintln(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log = log + fmt.Sprintln(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log = log + fmt.Sprintln(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log = log + fmt.Sprintln(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log = log + fmt.Sprintln(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log = log + fmt.Sprintln(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log = log + fmt.Sprintln(err)
	}

	err = w.Close()
	if err != nil {
		log = log + fmt.Sprintln(err)
	}

	c.Quit()

	return errors.New(log)
}
