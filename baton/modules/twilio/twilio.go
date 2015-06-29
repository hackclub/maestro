package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Twilio struct {
	UserId, ApiKey string
}

var URL = os.Getenv("HOSTLOCATION")

var (
	smsCallbacks  = make([]callback, 0)
	outboundCalls = make([]callback, 0)
	inboundCalls  = make(map[string]callback)
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
	case "send-mms":
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
	form := url.Values{"To": {to}, "From": {from}}

	if message, ok := body["body"]; ok {
		form.Add("Body", message.(string))
	}
	if url, ok := body["url"]; ok {
		form.Add("MediaUrl", url.(string))
	}

	res, err := t.postForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.UserId), form)
	if err != nil {
		log.Println("Twilio: Error in POST to /Messages.json")
		return err
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		log.Println("Twilio: Error decoding body as JSON")
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

	form := url.Values{"To": {to}, "From": {from}, "Url": {URL + "/baton/webhooks/Twilio/call/outbound"}}
	res, err := t.postForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", t.UserId), form)
	if err != nil {
		log.Println("Twilio: Error in POST to /Calls")
		return err
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		log.Println("Twilio: Error decoding body as JSON")
		return err
	}
	delete(jsonResponse, "account_sid")
	resp <- jsonResponse
	if message, ok := jsonResponse["message"]; ok {
		log.Println("Twilio: Error from Twilio server")
		return errors.New(message.(string))
	}

	outboundCalls = append(outboundCalls, callback{jsonResponse["to"].(string), resp, twiml})
	return nil
}

func (t Twilio) recieveSMS(body map[string]interface{}, resp chan<- interface{}) error {
	from := body["to"].(string)
	smsCallbacks = append(smsCallbacks, callback{from, resp, ""})
	return nil
}

func (t Twilio) recieveCall(body map[string]interface{}, resp chan<- interface{}) error {
	to := body["to"].(string)
	twiml := body["twiml"].(string)
	inboundCalls[to] = callback{to, resp, twiml}
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
	m.Path("/call/inbound").HandlerFunc(inboundCall)
	m.Path("/call/outbound").HandlerFunc(outboundCall)
	return m
}

func sms(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Twilio: Error parsing form")
		fmt.Println(err)
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	for i, callback := range smsCallbacks {
		if callback.number == jsonResponse["To"] {
			if !safeSend(callback.resp, jsonResponse) {
				smsCallbacks = append(smsCallbacks[:i], smsCallbacks[i+1:]...)
			}
		}
	}
}

func outboundCall(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Twilio: Error parsing form")
		log.Println(err)
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	for i, callback := range outboundCalls {
		if callback.number == jsonResponse["Called"] {
			fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
			safeSend(callback.resp, jsonResponse)
			outboundCalls = append(outboundCalls[:i], outboundCalls[i+1:]...)
			break
		}
	}
}

func inboundCall(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Twilio: Error parsing form")
		log.Println(err)
	}

	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	if callback, ok := inboundCalls[jsonResponse["Called"]]; ok {
		fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
		if !safeSend(callback.resp, jsonResponse) {
			delete(inboundCalls, jsonResponse["Called"])
		}
	}
}

//TODO: better system for safety
func safeSend(resp chan<- interface{}, data interface{}) bool {
	defer func() { recover() }()
	resp <- data
	return true
}
