package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type Twilio struct {
	UserId, ApiKey string
	smsCallbacks   []callback
}

func NewTwilio(userId string, apiKey string) Twilio {
	return Twilio{userId, apiKey, make([]callback, 0)}
}

type callback struct {
	number string
	resp   chan<- interface{}
}

var client = &http.Client{}

func (t Twilio) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	newBody := body.(map[string]interface{})
	switch cmd {
	case "send-sms":
		return t.sendSMS(newBody, resp)
	case "recieve-sms":
		return t.recieveSMS(newBody, resp)
	case "send-call":
		return nil
	case "recieve-call":
		return nil
	default:
		return errors.New("unknown command: " + cmd)
	}
}

func (t Twilio) sendSMS(body map[string]interface{}, resp chan<- interface{}) error {
	to := body["to"].(string)
	from := body["from"].(string)
	message := body["body"].(string)
	form := url.Values{"To": {to}, "From": {from}, "Body": {message}}

	res, err := t.postForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.UserId), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//Todo: filter what exports, currently it gives the account_sid probably not a good idea
	var out interface{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil
	}
	resp <- out
	return nil
}
func (t Twilio) recieveSMS(body map[string]interface{}, resp chan<- interface{}) error {
	from := body["from"].(string)
	t.smsCallbacks = append(t.smsCallbacks, callback{from, resp})
	return nil
}

func (t Twilio) postForm(url string, form url.Values) (*http.Response, error) {
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

	m.Path("/sms").HandlerFunc(makeSMS(t))
	return m
}
func makeSMS(t Twilio) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { t.sms(w, r) }
}
func (t Twilio) sms(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	out := make(map[string]string)
	for name, val := range r.PostForm {
		out[name] = val[0]
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		fmt.Println("Error Marshaling Form")
		return
	}
	fmt.Println("callbacks")
	for _, callback := range t.smsCallbacks {
		fmt.Println(callback, out["from"])
		if callback.number == out["from"] {
			fmt.Println("callback")
			callback.resp <- out
		}
	}
	fmt.Println(string(bytes))
}
