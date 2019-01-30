package main

import (
	"fmt"
	"github.com/gregdel/pushover"
	"regexp"
	"strings"
	"time"
)

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