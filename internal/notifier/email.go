package notifier

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

type emailSender interface {
	Send(ctx context.Context, subject, body string) error
}

type Client struct {
	host string
	port int
	user string
	pass string
	from string
	to   []string
}

func New(host, portStr, user, pass, from, toList string) *Client {
	if host == "" || portStr == "" || user == "" || pass == "" || from == "" || toList == "" {
		log.Println("[email] missing email config")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("[email] invalid EMAIL_PORT")
		return nil
	}
	return &Client{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
		to:   strings.Split(toList, ","),
	}
}

// Send отвправляет уведомление по email
func (c *Client) Send(ctx context.Context, subject, body string) error {
	if c == nil {
		return fmt.Errorf("[email] client is nil")
	}
	if len(c.to) == 0 {
		return fmt.Errorf("[email] no recipients in email")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	auth := smtp.PlainAuth("", c.user, c.pass, c.host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
			"\r\n%s",
		c.from, strings.Join(c.to, ","), subject, body,
	))

	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	err := smtp.SendMail(addr, auth, c.from, c.to, msg)
	if err != nil {
		log.Printf("[email] send error: %v", err)
		return err
	}
	log.Printf("[email] message successully sent to: %s", strings.Join(c.to, ","))
	return nil
}
