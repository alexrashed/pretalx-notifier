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

	// strip the message if it is longer than 1k characters (limit of pushover)
	if len(message) > 1024 {
		message = message[:1020] + "..."
	}

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