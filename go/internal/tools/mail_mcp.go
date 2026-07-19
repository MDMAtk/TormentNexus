package tools

import (
	"context"
	"fmt"
	"net/smtp"
)

func HandleSendEmail(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	to, _ :=getString(args, "to")
	if to == "" {
		return err("missing required argument: to")
}

	subject, _ :=getString(args, "subject")
	if subject == "" {
		return err("missing required argument: subject")
}

	body, _ :=getString(args, "body")
	if body == "" {
		return err("missing required argument: body")
}

	from, _ :=getString(args, "from")
	smtpHost, _ :=getString(args, "smtp_host")
	if smtpHost == "" {
		return err("missing required argument: smtp_host")
}

	smtpPort, _ :=getInt(args, "smtp_port")
	if smtpPort == 0 {
		return err("missing required argument: smtp_port")
}

	username, _ :=getString(args, "username")
	if username == "" {
		return err("missing required argument: username")
}

	password, _ :=getString(args, "password")
	if password == "" {
		return err("missing required argument: password")
}

	if from == "" {
		from = username
	}

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", username, password, smtpHost)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	sendErr := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if sendErr != nil {
		return err(sendErr.Error())
}

	return ok("Email sent successfully")
}

func HandleTestConnection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	smtpHost, _ :=getString(args, "smtp_host")
	if smtpHost == "" {
		return err("missing required argument: smtp_host")
}

	smtpPort, _ :=getInt(args, "smtp_port")
	if smtpPort == 0 {
		return err("missing required argument: smtp_port")
}

	username, _ :=getString(args, "username")
	if username == "" {
		return err("missing required argument: username")
}

	password, _ :=getString(args, "password")
	if password == "" {
		return err("missing required argument: password")
}

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", username, password, smtpHost)

	client, connErr := smtp.Dial(addr)
	if connErr != nil {
		return err(connErr.Error())
}

	defer client.Close()

	if authErr := client.Auth(auth); authErr != nil {
		return err(authErr.Error())
}

	return ok("Connection successful")
}