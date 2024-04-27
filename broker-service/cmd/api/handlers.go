package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action      string      `json:"action"`
	AuthPayload AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	c.writeJSON(w, http.StatusOK, payload)
}

func (c *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := c.readJSON(w, r, &requestPayload)
	if err != nil {
		c.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		c.authenticate(w, requestPayload.AuthPayload)
	default:
		c.errorJson(w, errors.New("unknown actions"))
	}
}

func (c *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		c.errorJson(w, errors.New("invalid Credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		c.errorJson(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		c.errorJson(w, err)
		return
	}

	if jsonFromService.Error {
		c.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	c.writeJSON(w, http.StatusAccepted, jsonFromService.Data)
}
