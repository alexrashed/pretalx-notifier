package main

import (
	"encoding/json"
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/jasonlvhit/gocron"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	//Image         interface{}   `json:"image"`
	//IsFeatured    bool          `json:"is_featured"`
	//Slot          interface{}   `json:"slot"`
	Speakers      []struct {
	//	Avatar    interface{} `json:"avatar"`
		Biography string      `json:"biography"`
		Code      string      `json:"code"`
		Name      string      `json:"name"`
	} `json:"speakers"`
	State          string `json:"state"`
	SubmissionType struct {
	//	De string `json:"de"`
		En string `json:"en"`
	} `json:"submission_type"`
	Title string      `json:"title"`
	//Track interface{} `json:"track"`
}

// submissions page json
type submissions struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous interface{}  `json:"previous"`
	Results  []submission `json:"results"`
}

// returns the string value of an environment variable or an error if not set or empty
func getEnvStr(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, fmt.Errorf(fmt.Sprintf("environment variable %s empty", key))
	}
	return value, nil
}

// returns the int value of an environment variable or an error if not set or empty
func getEnvInt(key string) (int, error) {
	str, err := getEnvStr(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// returns the boole value of an environment variable or an error if not set or empty
func getEnvBool(key string) (bool, error) {
	str, err := getEnvStr(key)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return value, nil
}

// reads the submissions from pretalx
func getSubmissions(pretalxURL string, pretalxAuthToken string) (submissions, error) {
	var submissions submissions
	// TODO implement paging instead of using a limit of 1k
	apiURL := pretalxURL + "/submissions?limit=1000"
	log.Print("requesting submissions from " + apiURL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return submissions,  errors.Wrap(err, "request could not be created")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+pretalxAuthToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return submissions, errors.Wrap(err, "request to PreTalx failed")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return submissions, errors.Wrap(err, "reading response body failed")
	}

	err = json.Unmarshal(body, &submissions)
	if err != nil {
		return submissions,  errors.Wrap(err, "un-marshalling response failed")
	}
	return submissions, nil
}

// sends the pushover notification on a new submission
func notify(pushOverApp *pushover.Pushover, pushoverUserToken, pretalxURL string, messages []string) (response *pushover.Response, err error) {
	title := fmt.Sprintf("PreTalx: Changed submissions detected")
	message := strings.Join(messages, "\n")
	recipient := pushover.NewRecipient(pushoverUserToken)

	// remove the API part from the URL
	r := regexp.MustCompile(`(.*).*/.*/.*/.+/?`)
	url := r.FindString(pretalxURL)

	messageObject := &pushover.Message{
		Message:     message,
		Title:       title,
		URL:         url,
		URLTitle:    "PreTalx",
		Timestamp:   time.Now().Unix(),
		Retry:       60 * time.Second,
		DeviceName:  "PreTalx-Notifier",
		Sound:       pushover.SoundCosmic,
	}
	response, err = pushOverApp.SendMessage(messageObject, recipient)
	return
}

// checks if there is any new submission on pretalx
func checkSubmissions(knownSubmissions map[string]submission, onlyNew bool, pretalxURL, pretalxAuthToken string,
	pushOverApp *pushover.Pushover, pushoverUserToken string) (err error) {
	messages := []string{}
	log.Print("checking for new submissions...")
	submissions, err := getSubmissions(pretalxURL, pretalxAuthToken)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, submission := range submissions.Results {
		if knownSubmission, found := knownSubmissions[submission.Code]; !found {
			log.Print("found a new submission, adding notification...")
			messages = append(messages, fmt.Sprintf("A new %s with the title '%s' has been submitted.", submission.SubmissionType.En, submission.Title))
			knownSubmissions[submission.Code] = submission
		} else if !onlyNew {
			if diff := pretty.Diff(knownSubmission, submission); diff != nil {
				log.Print("found a changed submission, adding notification...")
				messages = append(messages, fmt.Sprintf("The '%s' with the title '%s' has been changed: %s", submission.SubmissionType.En, submission.Title, diff))
				knownSubmissions[submission.Code] = submission
			}
		}
	}

	if len(messages) > 0 {
		response, err := notify(pushOverApp, pushoverUserToken, pretalxURL, messages)
		if err != nil {
			log.Fatal("sending notifications failed", err)
		} else {
			log.Print("notifications have been sent successfully:", response)
		}
	}

	log.Print("submission check completed")
	return
}

func main() {
	// check if all env variables are present
	pretalxURL, err := getEnvStr("PRETALX_URL")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pretalxAuthToken, err := getEnvStr("PRETALX_AUTH")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pushoverAPIToken, err := getEnvStr("PUSHOVER_API_TOKEN")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pushoverUserToken, err := getEnvStr("PUSHOVER_USER_TOKEN")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pushOverApp := pushover.New(pushoverAPIToken)
	minutes, err := getEnvInt("MINUTES")
	if err != nil {
		log.Print("MINUTES not set, using default of 15 minutes")
		minutes = 15
	}
	onlyNew, err := getEnvBool("ONLY_NEW")
	if err != nil {
		log.Print("ONLY_NEW not set, notifying on all changes")
		onlyNew = false
	}

	log.Print("initially downloading submissions...")
	submissions, err := getSubmissions(pretalxURL, pretalxAuthToken)
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

	// simply ignore failures and reschedule the task
	for ;; {
		gocron.Every(uint64(minutes)).Minutes().Do(checkSubmissions, knownSubmissions, onlyNew, pretalxURL, pretalxAuthToken, pushOverApp, pushoverUserToken)
		<-gocron.Start()
	}
}
