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
	data   string
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
		return t.makeCall(newBody, resp)
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

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		return err
	}
	delete(jsonResponse, "account_sid")
	resp <- jsonResponse
	return nil
}

func (t Twilio) makeCall(body map[string]interface{}, resp chan<- interface{}) error {
	to := body["to"].(string)
	from := body["from"].(string)
	twiml := body["twiml"].(string)
	form := url.Values{"To": {to}, "From": {from}, "Url": {"http://524b95fe.ngrok.io/baton/webhooks/Twilio/call"}} //temporary
	res, err := t.postForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", t.UserId), form)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		return err
	}
	delete(jsonResponse, "account_sid")
	if message, ok := jsonResponse["Message"]; ok {
		return errors.New(message.(string))
	}
	callCallbacks = append(callCallbacks, callback{to, resp, twiml})
	return nil
}
func (t Twilio) recieveSMS(body map[string]interface{}, resp chan<- interface{}) error {
	from := body["from"].(string)
	smsCallbacks = append(smsCallbacks, callback{from, resp, ""})
	return nil
}

func (t Twilio) recieveCall(body map[string]interface{}, resp chan<- interface{}) error {
	from := body["from"].(string)
	twiml := body["twiml"].(string)
	callCallbacks = append(callCallbacks, callback{from, resp, twiml})
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
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	for _, callback := range smsCallbacks {
		if callback.number == jsonResponse["From"] || callback.number == "*" {
			callback.resp <- jsonResponse
		}
	}

	fmt.Println(jsonResponse)
}
func call(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	for i, callback := range callCallbacks {
		if "inbound" == jsonResponse["Direction"] {
			if callback.number == jsonResponse["Caller"] {
				fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
				callback.resp <- jsonResponse
				callCallbacks = append(callCallbacks[:i], callCallbacks[i+1:]...)
				break
			}
		} else {
			if callback.number == jsonResponse["Called"] {
				fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
				callback.resp <- jsonResponse
				callCallbacks = append(callCallbacks[:i], callCallbacks[i+1:]...)
				break
			}
		}
	}

	fmt.Println(callCallbacks)
}
