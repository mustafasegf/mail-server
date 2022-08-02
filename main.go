package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Email    string
	Password string
  SMTPHost string
  SMTPPort string
}

type Request struct {
	Name    string
	Email   string
	Message string
}

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Panic("can't check path", err.Error())
	}
	file, err := os.OpenFile(path+"/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Panic("open file path", err.Error())
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}

func main() {
	config := Config{
		Email:    os.Getenv("EMAIL"),
		Password: os.Getenv("PASSWORD"),
    SMTPHost: os.Getenv("SMTP_HOST"),
    SMTPPort: os.Getenv("SMTP_PORT"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		to := []string{config.Email}

		body := "From: " + config.Email + "\n" +
			"To: " + strings.Join(to, ",") + "\n" +
			"Subject: " + "[SERVER] from: " + req.Email + "\n\n" +
			req.Message

		auth := smtp.PlainAuth("", config.Email, config.Password, config.SMTPHost)

		err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.Email, to, []byte(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		log.Println("Email sent", body)
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Listening on port", port)
	http.ListenAndServe(":"+port, nil)
}
