package giphy

import (
  "errors"
  "log"
  "bytes"
  "encoding/json"
  "github.com/gorilla/http"
)
type Giphy struct {
}

func (g Giphy) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	var url string
	switch cmd {
	  case "search":
	  	query := body.(map[string]interface{})["q"].(string)
	    url = "http://api.giphy.com/v1/gifs/search?api_key=dc6zaTOxFJmzC&q="+query
	  case "getbyid":
	    query := body.(map[string]interface{})["id"].(string)
	    url = "http://api.giphy.com/v1/gifs/"+query+"?api_key=dc6zaTOxFJmzC"
	  case "getbyids":
	    query := body.(map[string]interface{})["ids"].(string)
	    url = "http://api.giphy.com/v1/gifs?api_key=dc6zaTOxFJmzC&ids="+query
	  case "translate":
	    query := body.(map[string]interface{})["term"].(string)
	    url = "http://api.giphy.com/v1/gifs/translate?api_key=dc6zaTOxFJmzC&s="+query
	  case "random":
	    query := body.(map[string]interface{})["tags"].(string)
	    url = "http://api.giphy.com/v1/gifs/random?api_key=dc6zaTOxFJmzC&tags="+query
	  case "trending":
	    url = "http://api.giphy.com/v1/gifs/trending?api_key=dc6zaTOxFJmzC"
	  default:
	    return errors.New("unknown command: " + cmd)
	}
	log.Println(url)
	var tmp bytes.Buffer
  if _, err := http.Get(&tmp, url); err != nil {
      log.Fatalf("could not fetch: %v", err)
      return nil
  }
  var out interface{}
	if err := json.Unmarshal(tmp.Bytes(), &out); err != nil {
		log.Println("nu")
		log.Println(err)
		return nil
	}
	resp <- out
	return nil
}