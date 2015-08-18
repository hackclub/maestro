package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton"
)

var log = logrus.WithField("module", "Twilio")

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
	id     baton.CommandID
	data   string
}

var client = &http.Client{}

var response chan<- baton.Command

func (t Twilio) Init(cmd <-chan baton.Command, resp chan<- baton.Command) {
	response = resp
	go func() {
		for {
			go t.RunCommand(<-cmd)
		}
	}()
}

func (t Twilio) RunCommand(cmd baton.Command) {
	newBody := cmd.Body.(map[string]interface{})
	switch cmd.Call {
	case "send-sms":
		t.sendSMS(newBody, cmd.ID)
	case "send-mms":
		t.sendSMS(newBody, cmd.ID)
	case "recieve-sms":
		t.recieveSMS(newBody, cmd.ID)
	case "send-call":
		t.makeCall(newBody, cmd.ID)
	case "recieve-call":
		t.recieveCall(newBody, cmd.ID)
	default:
		log.WithField("command", cmd).Error("Unknown command")
	}
}

func (t Twilio) sendSMS(body map[string]interface{}, id baton.CommandID) {
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
		log.WithField("error", err).Error("Error in POST to /Messages.json")
		return
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		log.WithField("error", err).Error("Error decoding body as JSON")
		return
	}
	delete(jsonResponse, "account_sid")
	send(id, "send-sms", jsonResponse)
}

func (t Twilio) makeCall(body map[string]interface{}, id baton.CommandID) {
	to := body["to"].(string)
	from := body["from"].(string)
	twiml := body["twiml"].(string)

	form := url.Values{"To": {to}, "From": {from}, "Url": {URL + "/baton/webhooks/Twilio/call/outbound"}}
	res, err := t.postForm(fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", t.UserId), form)
	if err != nil {
		log.WithField("error", err).Error("Error in POST to /Calls")
		return
	}
	defer res.Body.Close()

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		log.WithField("error", err).Error("Error decoding body as JSON")
		return
	}
	delete(jsonResponse, "account_sid")
	send(id, "send-call", jsonResponse)
	if _, ok := jsonResponse["message"]; ok {
		log.WithField("error", jsonResponse).Error("Error from Twilio server")
		return
	}
	outboundCalls = append(outboundCalls, callback{jsonResponse["to"].(string), id, twiml})
}

func (t Twilio) recieveSMS(body map[string]interface{}, id baton.CommandID) {
	from := body["to"].(string)
	smsCallbacks = append(smsCallbacks, callback{from, id, ""})
}

func (t Twilio) recieveCall(body map[string]interface{}, id baton.CommandID) {
	to := body["to"].(string)
	twiml := body["twiml"].(string)
	inboundCalls[to] = callback{to, id, twiml}
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
		log.WithField("error", err).Error("Error parsing form")
		return
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	log.WithField("number", jsonResponse["To"]).Info("SMS recieved on number")
	for i, callback := range smsCallbacks {
		if callback.number == jsonResponse["To"] {
			log.WithFields(logrus.Fields{
				"number": jsonResponse["To"],
				"id":     callback.id,
			}).Debug("Callback sent to")
			if !send(callback.id, "recieve-sms", jsonResponse) {
				smsCallbacks = append(smsCallbacks[:i], smsCallbacks[i+1:]...)
			}
		}
	}
}

func outboundCall(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.WithField("error", err).Error("Error parsing form")
		return
	}
	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	log.WithField("number", jsonResponse["Called"]).Info("Outbound call")
	for i, callback := range outboundCalls {
		if callback.number == jsonResponse["Called"] {
			fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
			send(callback.id, "send-call", jsonResponse)
			outboundCalls = append(outboundCalls[:i], outboundCalls[i+1:]...)
			break
		}
	}
}

func inboundCall(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.WithField("error", err).Error("Error parsing form")
		return
	}

	jsonResponse := make(map[string]string)
	for name, val := range r.PostForm {
		jsonResponse[name] = val[0]
	}
	delete(jsonResponse, "AccountSid")
	log.WithField("number", jsonResponse["Called"]).Info("Inbound call")
	if callback, ok := inboundCalls[jsonResponse["Called"]]; ok {
		fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?><Response>%s</Response>", callback.data)
		if !send(callback.id, "recieve-call", jsonResponse) {
			delete(inboundCalls, jsonResponse["Called"])
		}
	}
}

func send(id baton.CommandID, call string, body interface{}) bool {
	response <- baton.Command{"Twilio", call, id, body}
	return true //will probably be replaced later
}
