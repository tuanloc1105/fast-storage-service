package implement

import (
	"context"
	"crypto/tls"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailProperties struct {
	To                 []string
	Cc                 []string
	Bcc                []string
	Subject            string
	Content            string
	PathOfFileToEmbed  []string
	PathOfFileToAttach []string
}

func sendHtmlEmailContent(ctx context.Context, properties EmailProperties) error {
	// smtpHost := "smtp.office365.com"
	// smtpPort := 587
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("OUTLOOK_USERNAME")
	password := os.Getenv("OUTLOOK_PASSWORD")

	if smtpHost == "" || smtpPort == "" {
		return errors.New("`SMTP_HOST` and `SMTP_PORT` must be set in the environment variable")
	}

	if username == "" || password == "" {
		return errors.New("`OUTLOOK_USERNAME` and `OUTLOOK_PASSWORD` must be set in the environment variable")
	}

	if properties.To == nil || len(properties.To) == 0 {
		return errors.New("email recipient cannot be left blank")
	}

	if properties.Subject == "" {
		return errors.New("email subject cannot be left blank")
	}

	if properties.Content == "" {
		return errors.New("email content cannot be left blank")
	}

	properties.Content = strings.Trim(properties.Content, "\n")

	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("From", username)
	goMailMessage.SetHeader("To", properties.To...)
	goMailMessage.SetHeader("Cc", properties.Cc...)
	goMailMessage.SetHeader("Bcc", properties.Bcc...)
	goMailMessage.SetHeader("Subject", properties.Subject)
	goMailMessage.SetBody("text/html", properties.Content)
	if len(properties.PathOfFileToEmbed) > 0 {
		for indexOfFileToBeEmbeded, fileToBeEmbeded := range properties.PathOfFileToEmbed {
			goMailMessage.Embed(fileToBeEmbeded, gomail.Rename(fmt.Sprintf("image-%d", indexOfFileToBeEmbeded)))
		}
	}
	if len(properties.PathOfFileToAttach) > 0 {
		for _, fileToBeAttached := range properties.PathOfFileToAttach {
			goMailMessage.Attach(fileToBeAttached)
		}
	}

	smtpPortInt, smtpPortConvertToIntError := strconv.Atoi(smtpPort)

	if smtpPortConvertToIntError != nil {
		log.WithLevel(constant.Error, ctx, "Could not send email: %v", smtpPortConvertToIntError)
		return smtpPortConvertToIntError
	}

	goMailDialer := gomail.NewDialer(smtpHost, smtpPortInt, username, password)
	goMailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if sendEmailError := goMailDialer.DialAndSend(goMailMessage); sendEmailError != nil {
		log.WithLevel(constant.Error, ctx, "Could not send email: %v", sendEmailError)
		return sendEmailError
	}

	log.WithLevel(constant.Info, ctx, "Email sent successfully")
	return nil
}
