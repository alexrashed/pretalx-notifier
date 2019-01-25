package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// submission json
type submission struct {
	Abstract      string        `json:"abstract"`
	Answers       []interface{} `json:"answers"`
	Code          string        `json:"code"`
	ContentLocale string        `json:"content_locale"`
	Description   string        `json:"description"`
	DoNotRecord   bool          `json:"do_not_record"`
	Duration      string        `json:"duration"`
	Image         interface{}   `json:"image"`
	IsFeatured    bool          `json:"is_featured"`
	Slot          interface{}   `json:"slot"`
	Speakers      []struct {
		Avatar    interface{} `json:"avatar"`
		Biography string      `json:"biography"`
		Code      string      `json:"code"`
		Name      string      `json:"name"`
	} `json:"speakers"`
	State          string `json:"state"`
	SubmissionType struct {
		De string `json:"de"`
		En string `json:"en"`
	} `json:"submission_type"`
	Title string      `json:"title"`
	Track interface{} `json:"track"`
}

// submissions page json
type submissions struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous interface{}  `json:"previous"`
	Results  []submission `json:"results"`
}

// reads the string value of an environment variable or an error if not set or empty
func getEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, errors.New(fmt.Sprintf("environment variable %s empty", key))
	}
	return v, nil
}

// reads the int value of an environment variable or an error if not set or empty
func getEnvInt(key string) (int, error) {
	s, err := getEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// reads the submissions from pretalx
func getSubmissions(pretalxUrl string, pretalxAuthToken string) (submissions, error) {
	var submissions submissions
	// TODO implement paging instead of using a limit of 1k
	apiUrl := pretalxUrl + "/submissions?limit=1000"
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

// sends the pushover notification on a new submission
func notify(pushOverApp *pushover.Pushover, pushoverUserToken string, submission submission) (response *pushover.Response, err error) {
	title := fmt.Sprintf("PreTalx: New %s", submission.SubmissionType.En)
	message := fmt.Sprintf("A new %s with the title '%s' has been submitted.", submission.SubmissionType.En, submission.Title)
	recipient := pushover.NewRecipient(pushoverUserToken)
	messageObject := pushover.NewMessageWithTitle(message, title)
	response, err = pushOverApp.SendMessage(messageObject, recipient)
	return
}

// checks if there is any new submission on pretalx
func checkSubmissions(knownSubmissions map[string]submission, pretalxUrl, pretalxAuthToken string,
	pushOverApp *pushover.Pushover, pushoverUserToken string) (err error) {
	log.Print("checking for new submissions...")
	submissions, err := getSubmissions(pretalxUrl, pretalxAuthToken)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, submission := range submissions.Results {
		if _, found := knownSubmissions[submission.Code]; !found {
			log.Print("found a new submission, sending notification...")
			// a new submission has been found
			response, err := notify(pushOverApp, pushoverUserToken, submission)
			if err != nil {
				log.Fatal("sending notification on new submission failed")
			} else {
				log.Print("notification has been sent successfully:")
				log.Print(response)
			}
			knownSubmissions[submission.Code] = submission
		}
	}
	log.Print("submission check completed")
	return
}

func main() {
	// check if all env variables are present
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
	pushoverApiToken, err := getEnvStr("PUSHOVER_API_TOKEN")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pushoverUserToken, err := getEnvStr("PUSHOVER_USER_TOKEN")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pushOverApp := pushover.New(pushoverApiToken)
	minutes, err := getEnvInt("MINUTES")
	if err != nil {
		log.Print("no minutes set, using default of 15 minutes")
		minutes = 15
	}

	log.Print("initially downloading submissions...")
	submissions, err := getSubmissions(pretalxUrl, pretalxAuthToken)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	var knownSubmissions = make(map[string]submission)

	// initially download the submissiosn (without notifying)
	for _, submission := range submissions.Results {
		knownSubmissions[submission.Code] = submission
	}

	// schedule checkSubmissions every x minutes
	log.Print(fmt.Sprintf("scheduling submission check every %d minutes", minutes))
	gocron.Every(uint64(minutes)).Minutes().Do(checkSubmissions, knownSubmissions, pretalxUrl, pretalxAuthToken, pushOverApp, pushoverUserToken)
	<- gocron.Start()
}