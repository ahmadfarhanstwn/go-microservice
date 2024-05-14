package main

import (
	"log"
	"net/http"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (c *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload mailMessage

	err := c.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Print(err)
		c.errorJson(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = c.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Print(err)
		c.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	c.writeJSON(w, http.StatusAccepted, payload)
}
