package main

import (
	"fmt"
	"github.com/gregdel/pushover"
	"github.com/jasonlvhit/gocron"
	"log"
)

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

	// initially download the submission (without notifying)
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
