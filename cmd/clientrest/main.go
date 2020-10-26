package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	address := "http://localhost:8080"
	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	// call create
	resp, err := http.Post(
		address+"/v1/todo",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
     		"toDo": {
				"title":"title (%s)",
				"description":"dascription (%s)",
				"reminder":"%s"
			}
		}`, pfx, pfx, pfx)),
	)
	if err != nil {
		log.Fatalf("create error: %+v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("create response body read error: %+v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	var created struct {
		API string `json:"api"`
		ID  string `json:"id"`
	}

	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("json response of create method error: %+v", err)
	}

	// call read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s?api=v1", address, "/v1/todo", created.ID))
	if err != nil {
		log.Fatalf("read error: %+v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("read response body error: %+v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// call readall
	resp, err = http.Get(fmt.Sprintf("%s%s/all?api=v1", address, "/v1/todo"))
	if err != nil {
		log.Fatalf("read all error: %+v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("read all response body error: %+v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("read all response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	//call update
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s%s/%s", address, "/v1/todo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
     		"toDo": {
				"title":"title (%s) update",
				"description":"dascription (%s) update",
				"reminder":"%s"
			}
		}`, pfx, pfx, pfx)),
	)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("update method error: %+v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("update response body error: %+v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// call delete
	req, err = http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s%s/%s?api=v1", address, "/v1/todo", created.ID),
		nil,
	)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call delete method: %+v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("delete response body error: %+v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
