package src

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Category struct {
  Name string `json:"name"`;
  Description string `json:"description"`;
  Color string `json:"color"`;
  BgColor string `json:"bgColor"`;
  FgColor string `json:"fgColor"`
}

func GetCategory(c *http.Client, url string, name string) {
  body := new(bytes.Buffer)
  request, requestError := http.NewRequest("GET", url, body)
  if requestError != nil {
    log.Println(requestError.Error())
  } 
  res, resErr := c.Do(request)
  if resErr != nil {
    log.Println(resErr)
  }
  defer res.Body.Close()
  var resData interface{}
  json.NewDecoder(res.Body).Decode(&resData)
  log.Println(resData)
  log.Println("no")
}
