package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func GetCategory(c *http.Client, url string) (Category, error) {
	body := new(bytes.Buffer)
	request, requestError := http.NewRequest("GET", url, body)
	if requestError != nil {
		log.Println(requestError.Error())
		return Category{}, requestError
	}
	res, resErr := c.Do(request)
	if resErr != nil {
		log.Println(resErr)
		return Category{}, resErr
	}
	defer res.Body.Close()
	if status := res.StatusCode; status != http.StatusOK {
		var resErr HttpResponseError
		err := json.NewDecoder(res.Body).Decode(&resErr)
		if err != nil {
			return Category{}, err
		}
		log.Println(resErr.Message)
		return Category{}, errors.New(resErr.Message)
	}
	var resData HttpResponse[Category]
	err := json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return Category{}, err
	}
	return resData.Data, nil
}

func CreateCategory(c *http.Client, url string, category Category) (Category, error) {
	const DefaultRbgaColor = "rgba(0,0,0,1)"
	category.FgColor = DefaultRbgaColor
	category.BgColor = DefaultRbgaColor
	category.Color = DefaultRbgaColor
	categoryJson, jsonErr := json.Marshal(category)
	if jsonErr != nil {
		log.Println("category.go 48:  " + jsonErr.Error())
		return Category{}, jsonErr
	}
	body := bytes.NewBuffer(categoryJson)
	request, requestError := http.NewRequest("POST", url, body)
	if requestError != nil {
		log.Println("category.go 54:  " + requestError.Error())
		return Category{}, requestError
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	res, resErr := c.Do(request)
	if resErr != nil {
		log.Println("category.go 63:  " + resErr.Error())
		return Category{}, resErr
	}
	defer res.Body.Close()
	if status := res.StatusCode; status == http.StatusNotFound {
		return Category{}, errors.New("category create endpoint not found")
	}
	if status := res.StatusCode; status >= http.StatusInternalServerError {
		var resErr SpringBootResponseError
		err := json.NewDecoder(res.Body).Decode(&resErr)
		if err != nil {
			log.Println("category.go 67:  " + err.Error())
			return Category{}, err
		}
		log.Println("category.go 68: ", resErr)
		return Category{}, errors.New(resErr.Message)
	}
	if status := res.StatusCode; status != http.StatusCreated {
		var resErr HttpResponseError
		err := json.NewDecoder(res.Body).Decode(&resErr)
		if err != nil {
			log.Println("category.go 69:  " + err.Error())
			return Category{}, err
		}
		return Category{}, errors.New(resErr.Message)
	}
	var resData HttpResponse[Category]
	err := json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Println("category.go 79:  " + err.Error())
		return Category{}, err
	}
	return resData.Data, nil
}
