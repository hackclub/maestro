package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type Twilio struct {
	UserId, ApiKey string
}

var client = &http.Client{}

func (t Twilio) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	newBody := body.(map[string]interface{})
	to := newBody["to"].(string)
	from := newBody["from"].(string)
	message := newBody["body"].(string)
	form := url.Values{"To": {to}, "From": {from}, "Body": {message}}

	_, err := t.PostForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.UserId), form)
	if err != nil {
		return err
	}

	return nil
}

func (t Twilio) PostForm(url string, form url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.UserId, t.ApiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t Twilio) Handler() *mux.Router {
	m := mux.NewRouter()
	return m
}
