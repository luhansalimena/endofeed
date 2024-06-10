package main

import (
	"encoding/json"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/imroc/req/v3"
)

type announcements struct {
	Data []announcement `json:"announcements"`
}

type announcement struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func main() {
	var announcements []announcement
	cities := []string{"niterói", "rio de janeiro", "são gonçalo"}
	for _, city := range cities {
		data := getAnnouncement(city)
		announcements = append(announcements, data.Data...)
	}
	for _, announcement := range announcements {
		date, _ := time.Parse(time.RFC3339, announcement.CreatedAt)
		if date.Day() == 10 && date.Month() == 4 {
			sendEmail()
		}
	}
}

func getAnnouncement(city string) announcements {
	client := req.C()
	resp, err := client.R().SetQueryParam("city", city).Get("https://classificados-api.cro-rj.org.br/announcements?district=&category=4&page=1")
	if err != nil {
		log.Fatal(err)
	}
	body, _ := resp.ToString()
	var announcements announcements
	_ = json.Unmarshal([]byte(body), &announcements)

	return announcements
}

func sendEmail() {
	username := os.Getenv("MAIL_SMTP_USERNAME")
	password := os.Getenv("MAIL_SMTP_PASSWORD")
	smtpHost := os.Getenv("MAIL_SMTP_HOST")

	// Choose auth method and set it up

	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Message data

	from := os.Getenv("MAIL_FROM")

	to := []string{os.Getenv("MAIL_TO")}

	message := []byte(
		"Subject: Novas Oportunidades\r\n" +
			"\r\n" +
			"Olá, você recebeu uma nova oportunidade de trabalho")

	smtpUrl := smtpHost + ":2525"

	err := smtp.SendMail(smtpUrl, auth, from, to, message)

	if err != nil {
		log.Fatal(err)
	}

}
