package mails

import (
	"log"
	"os"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

var (
	host     = "smtp.gmail.com"
	port     = "587"
	from     = os.Getenv("MAIL")
	password = os.Getenv("MAIL_PASS")
)

func SendExpertConfirmationNotification(toEmail string) {
	auth := sasl.NewPlainClient("", from, password)

	to := []string{toEmail}
	msg := strings.NewReader("To: " + toEmail + "\r\n" +
		"Subject: You are now confirmed expert in Mindflow\r\n" +
		"\r\n" +
		"You are now confirmed as an expert in Mindflow")

	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		log.Println(err)
	}
}

func SendExpertRejectNotification(toEmail string) {
	auth := sasl.NewPlainClient("", from, password)

	to := []string{toEmail}
	msg := strings.NewReader("To: " + toEmail + "\r\n" +
		"Subject: You were rejected to become expert in Mindflow\r\n" +
		"\r\n" +
		"You were rejected to become expert in Mindflow")

	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		log.Println(err)
	}
}

func SendConsultationNotification(toEmails []string, link string) {
	auth := sasl.NewPlainClient("", from, password)

	for _, to := range toEmails {
		msg := strings.NewReader("To: " + to + "\r\n" +
			"Subject: You were invived to consultation\r\n" +
			"\r\n" +
			"You were invived to consultation\r\n" +
			"Link: " + link)
		toSlice := []string{to}
		err := smtp.SendMail(host+":"+port, auth, from, toSlice, msg)
		if err != nil {
			log.Println(err)
		}
	}
}
