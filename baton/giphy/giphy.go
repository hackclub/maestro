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

func (e Giphy) RunCommand(cmd string, body interface{}, resp chan<- interface{}) error {
	if cmd != "search" {
		return errors.New("unknown command: " + cmd)
	}
	query := body.(map[string]interface{})["q"].(string)
	log.Println(query);
	var tmp bytes.Buffer
  if _, err := http.Get(&tmp, "http://api.giphy.com/v1/gifs/search?q="+query+"&api_key=dc6zaTOxFJmzC"); err != nil {
      log.Fatalf("could not fetch: %v", err)
      return nil
  }
  log.Println(tmp.String())
  var out interface{}
	if err := json.Unmarshal(tmp.Bytes(), &out); err != nil {
		log.Println("nu")
		log.Println(err)
		return nil
	}
	resp <- out
	return nil
}