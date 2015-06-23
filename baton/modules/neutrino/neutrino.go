package neutrino

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Neutrino struct {
	UserId, ApiKey string
}

func (n Neutrino) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	v := url.Values{}
	newBody := body.(map[string]interface{})
	switch cmd {
	case "geocode-address":
		v.Add("address", newBody["address"].(string))
	case "geocode-reverse":
		v.Add("latitude", newBody["latitude"].(string))
		v.Add("longitude", newBody["longitude"].(string))
	case "qr-code":
	case "bin-lookup":
		v.Add("bin-number", newBody["bin-number"].(string))
	case "bad-word-filter":
		v.Add("content", newBody["content"].(string))
	case "convert":
		v.Add("from-value", newBody["from-value"].(string))
		v.Add("from-type", newBody["from-type"].(string))
		v.Add("to-type", newBody["to-type"].(string))
	default:
		return errors.New("unknown command: " + cmd)
	}
	v.Add("user-id", n.UserId)
	v.Add("api-key", n.ApiKey)
	v.Add("ip", "162.209.104.195")
	url := "https://neutrinoapi.com/" + cmd

	data, err := http.PostForm(url, v)
	if err != nil {
		log.Println("Neutrino: Could not POST data")
		return err
	}
	defer data.Body.Close()

	var out interface{}
	if err := json.NewDecoder(data.Body).Decode(&out); err != nil {
		log.Println("Neutrino: Error decoding body as JSON")
		return err
	}
	resp <- out
	return nil
}

func (n Neutrino) Handler() *mux.Router {
	return mux.NewRouter()
}
