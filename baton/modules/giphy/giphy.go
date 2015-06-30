package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/hackedu/maestro/baton/commands"
)

type Giphy struct {
	ApiKey string
}

var resp chan<- commands.Command

func (g Giphy) Init(cmd <-chan commands.Command, resp chan<- commands.Command) {
	resp = resp
	go func() {
		for {
			go g.RunCommand(<-cmd)
		}
	}()
}

func (g Giphy) RunCommand(cmd commands.Command) error {
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
		return errors.New("unknown command: " + cmd.Call)
	}
	if err != nil {
		log.Println("Giphy: error making URL")
		return err
	}
	log.Println("Giphy: URL to be requested", u.String())
	res, err := http.Get(u.String())
	if err != nil {
		log.Println("Giphy: Error in GET request")
		return err
	}
	defer res.Body.Close()

	var out interface{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		log.Println("Giphy: Error decoding body as JSON")
		return err
	}
	resp <- commands.Command{"Giphy", cmd.Call, cmd.ID, out}
	return nil
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
func (g Giphy) Handler() *mux.Router {
	return mux.NewRouter()
}
