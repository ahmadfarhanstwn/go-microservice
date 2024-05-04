package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.readJSON(w, r, &requestPayload)
	if err != nil {
		c.errorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := c.models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		c.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		c.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// send log
	err = c.sendLog("authenticate", "request authenticate")
	if err != nil || !valid {
		c.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	c.writeJSON(w, http.StatusAccepted, payload)
}

func (c *Config) sendLog(name string, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
