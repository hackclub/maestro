package neutrino

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton/commands"
)

type Neutrino struct {
	UserId, ApiKey string
}

var resp chan<- commands.Command

func (n Neutrino) Init(cmd <-chan commands.Command, resp chan<- commands.Command) {
	resp = resp
	go func() {
		for {
			go n.RunCommand(<-cmd)
		}
	}()
}

func (n Neutrino) RunCommand(cmd commands.Command) {
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

	data, err := http.PostForm(url, v)
	if err != nil {
		log.Println("Neutrino: Could not POST data")
		log.Println(err)
	}
	defer data.Body.Close()

	var out interface{}
	if err := json.NewDecoder(data.Body).Decode(&out); err != nil {
		log.Println("Neutrino: Error decoding body as JSON")
		log.Println(err)
	}
	resp <- commands.Command{"Neutrino", cmd.Call, cmd.ID, out}
}

func (n Neutrino) Handler() *mux.Router {
	return mux.NewRouter()
}
