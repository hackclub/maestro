package giphy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/hackedu/maestro/baton"
)

var log = logrus.WithField("module", "Giphy")

type Giphy struct {
	ApiKey string
}

var resp chan<- baton.Command

func (g Giphy) Init(cmd <-chan baton.Command, resp chan<- baton.Command) {
	go func() {
		for {
			go g.RunCommand(<-cmd, resp)
		}
	}()
}

func (g Giphy) RunCommand(cmd baton.Command, resp chan<- baton.Command) {
	var u url.URL
	var err error
	switch cmd.Call {
	case "search":
		query := cmd.Body.(map[string]interface{})["q"].(string)
		u, err = g.makeURL("gifs/search", url.Values{"q": {query}})
	case "getbyid":
		id := cmd.Body.(map[string]interface{})["id"].(string)
		u, err = g.makeURL(fmt.Sprintf("gifs/%s", id), url.Values{})
	case "getbyids":
		ids := cmd.Body.(map[string]interface{})["ids"].(string)
		u, err = g.makeURL("gifs", url.Values{"ids": {ids}})
	case "translate":
		term := cmd.Body.(map[string]interface{})["term"].(string)
		u, err = g.makeURL("gifs/translate", url.Values{"s": {term}})
	case "random":
		tags := cmd.Body.(map[string]interface{})["tags"].(string)
		u, err = g.makeURL("gifs/random", url.Values{"tags": {tags}})
	case "trending":
		u, err = g.makeURL("gifs/trending", url.Values{})
	default:
		log.WithField("command", cmd).Error("Unknown Command")
		return
	}
	if err != nil {
		log.WithField("error", err).Error("Error making URL")
		return
	}

	log.WithField("URL", u.String()).Debug("URL to be requested")
	res, err := http.Get(u.String())
	if err != nil {
		log.WithField("error", err).Error("Error in GET request")
		return
	}
	defer res.Body.Close()

	var out interface{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		log.WithField("error", err).Error("Error decoding body as JSON")
		return
	}
	resp <- baton.Command{"Giphy", cmd.Call, cmd.ID, out}
}

func (g Giphy) makeURL(path string, v url.Values) (url.URL, error) {
	baseURL, _ := url.Parse("http://api.giphy.com/v1/")
	v.Add("api_key", g.ApiKey)

	rel, err := url.Parse(path)
	if err != nil {
		return url.URL{}, err
	}

	u := baseURL.ResolveReference(rel)
	u.RawQuery = v.Encode()
	return *u, nil
}
