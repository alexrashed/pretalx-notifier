package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type submissions struct {
	count int
}

func getEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, errors.New(fmt.Sprintf("environment variable %s empty", key))
	}
	return v, nil
}

func getSubmissions(pretalxUrl string, pretalxAuthToken string) (submissions, error) {
	var submissions submissions
	apiUrl := pretalxUrl + "/submissions"
	log.Print("requesting submissions from " + apiUrl)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return submissions, errors.New("request could not be created")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token " + pretalxAuthToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return submissions, errors.New("request to PreTalx failed")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return submissions, errors.New("reading response body failed")
	}

	err = json.Unmarshal(body, &submissions)
	if err != nil {
		return submissions, errors.New("un-marshalling response failed")
	}
	return submissions, nil
}

func main() {
	pretalxUrl, err := getEnvStr("PRETALX_URL")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pretalxAuthToken, err := getEnvStr("PRETALX_AUTH")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	submissions, err := getSubmissions(pretalxUrl, pretalxAuthToken)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Print(submissions)

	log.Print(pretalxUrl)
	log.Print(pretalxAuthToken)

	log.Print("retrieving all submissions...")
}