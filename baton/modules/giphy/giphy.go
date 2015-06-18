package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Giphy struct {
	ApiKey string
}

func (g Giphy) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	var u url.URL
	var err error
	switch cmd {
	case "search":
		query := body.(map[string]interface{})["q"].(string)
		u, err = g.makeURL("gifs/search", url.Values{"q": {query}})
	case "getbyid":
		id := body.(map[string]interface{})["id"].(string)
		u, err = g.makeURL(fmt.Sprintf("gifs/%s", id), url.Values{})
	case "getbyids":
		ids := body.(map[string]interface{})["ids"].(string)
		u, err = g.makeURL("gifs", url.Values{"ids": {ids}})
	case "translate":
		term := body.(map[string]interface{})["term"].(string)
		u, err = g.makeURL("gifs/translate", url.Values{"s": {term}})
	case "random":
		tags := body.(map[string]interface{})["tags"].(string)
		u, err = g.makeURL("gifs/random", url.Values{"tags": {tags}})
	case "trending":
		u, err = g.makeURL("gifs/trending", url.Values{})
	default:
		return errors.New("unknown command: " + cmd)
	}
	log.Println(u.String())
	if err != nil {
		return err
	}
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var out interface{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		log.Println("nu")
		log.Println(err)
		return nil
	}
	resp <- out
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
