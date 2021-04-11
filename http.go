package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func patchPage(payload []byte, url string, token string) {
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))
}
