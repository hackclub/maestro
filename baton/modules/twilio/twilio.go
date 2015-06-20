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
}

var (
	smsCallbacks  = make([]callback, 0)
	callCallbacks = make([]callback, 0)
)

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
		return t.recieveCall(newBody, resp)
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
	smsCallbacks = append(smsCallbacks, callback{from, resp})
	return nil
}

func (t Twilio) recieveCall(body map[string]interface{}, resp chan<- interface{}) error {
	from := body["from"].(string)
	callCallbacks = append(callCallbacks, callback{from, resp})
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

	m.Path("/sms").HandlerFunc(sms)
	m.Path("/call").HandlerFunc(call)
	return m
}
func sms(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	out := make(map[string]string)
	for name, val := range r.PostForm {
		out[name] = val[0]
	}
	for _, callback := range smsCallbacks {
		if callback.number == out["From"] {
			callback.resp <- out
		}
	}

	fmt.Println(out)
}
func call(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	out := make(map[string]string)
	for name, val := range r.PostForm {
		out[name] = val[0]
	}
	for _, callback := range callCallbacks {
		if "inbound" == out["Direction"] {
			if callback.number == out["Caller"] {
				callback.resp <- out
			}
		} else {
			if callback.number == out["Called"] {
				callback.resp <- out
			}
		}
	}

	fmt.Println(out)
}
