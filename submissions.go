package main

import (
	"encoding/json"
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

// checks if there is any new submission on pretalx
func checkSubmissions(knownSubmissions map[string]submission, onlyNew bool, pretalxURL, pretalxAuthToken string,
	pushOverApp *pushover.Pushover, pushoverUserToken string) (err error) {
	var messages []string
	log.Print("checking for new submissions...")
	submissions, err := getSubmissions(pretalxURL, pretalxAuthToken)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, submission := range submissions.Results {
		if knownSubmission, found := knownSubmissions[submission.Code]; !found {
			log.Print("found a new submission, adding notification...")
			messages = append(messages, fmt.Sprintf("A new %s with the title '%s' has been submitted.", getSubmissionTypeString(submission.SubmissionType), submission.Title))
			knownSubmissions[submission.Code] = submission
		} else if !onlyNew {
			if diff := pretty.Diff(knownSubmission, submission); diff != nil {
				log.Print("found a changed submission, adding notification...")
				messages = append(messages, fmt.Sprintf("The %s with the title '%s' has been changed: %s", getSubmissionTypeString(submission.SubmissionType), submission.Title, diff))
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

// get the name of the submission type (prefer english, use german as fallback)
func getSubmissionTypeString(st submissionType) string {
	if st.En != "" {
		return st.En
	} else {
		return st.De
	}
}
