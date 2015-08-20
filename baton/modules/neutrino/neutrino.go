package neutrino

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/hackedu/maestro/baton"
)

var log = logrus.WithField("module", "Neutrino")

type Neutrino struct {
	UserId, ApiKey string
}

func (n Neutrino) Init(cmd <-chan baton.Command, resp chan<- baton.Command) {
	go func() {
		for {
			go n.RunCommand(<-cmd, resp)
		}
	}()
}

func (n Neutrino) RunCommand(cmd baton.Command, resp chan<- baton.Command) {
	v := url.Values{}
	newBody := cmd.Body.(map[string]interface{})
	switch cmd.Call {
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
		log.Println("unknown command: " + cmd.Call)
	}
	v.Add("user-id", n.UserId)
	v.Add("api-key", n.ApiKey)
	v.Add("ip", "162.209.104.195")
	url := "https://neutrinoapi.com/" + cmd.Call

	log.WithField("URL", url).Debug("URL to be requested")

	data, err := http.PostForm(url, v)
	if err != nil {
		log.WithField("error", err).Error("Could not POST data")
		return
	}
	defer data.Body.Close()

	var out interface{}
	if err := json.NewDecoder(data.Body).Decode(&out); err != nil {
		log.WithField("error", err).Error("Error decoding body as JSON")
		return
	}
	resp <- baton.Command{"Neutrino", cmd.Call, cmd.ID, out}
}
