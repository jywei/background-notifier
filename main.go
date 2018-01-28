package main

import (
	"github.com/benmanns/goworker"
	"fmt"
	"time"
	"github.com/parnurzeal/gorequest"
)

type Payload struct {
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Text      string `json:"text,omitempty"`
}


func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}


func send(webhookURL string, proxy string, payload Payload) []error {
	request := gorequest.New().Proxy(proxy)
	resp, _, err := request.
		Post(webhookURL).
		RedirectPolicy(redirectPolicyFunc).
		Send(payload).
		End()

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return []error{fmt.Errorf("Error sending msg. Status: %v", resp.Status)}
	}
	return nil
}

func notificationFunc(payload Payload, webhookURL string, response chan<- []error)  {
	response <- send(webhookURL, "", payload)
}

func notificationWorker(queue string, args ...interface{}) error {
	webhookURL := "https://hooks.slack.com/services/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"
	responseChannel := make(chan []error)
	payload := Payload{
		Text:      args[0].(string),
		Username:  "Hexter Bot",
		IconEmoji: ":hexter_is_ur_daddy:",
		Channel:   "#nuclear-testing-sites",
	}
	go notificationFunc(payload, webhookURL, responseChannel)
	fmt.Printf("Send to %s Platform, message content: %v\n", queue, args)
	err := <- responseChannel
	return err[0]
}

func init() {
	settings := goworker.WorkerSettings{
		URI: "redis://localhost:6379/",
		Connections: 100,
		Queues:	[]string{"slack", "message"},
		UseNumber: true,
		ExitOnComplete: false,
		Concurrency: 2,
		Namespace: "resque:",
		Interval: 5.0,
	}
	goworker.SetSettings(settings)
	goworker.Register("notifier", notificationWorker)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error", err)
	}
	fmt.Printf("Started on %v", time.Now().Format("2006-01-02 15:04:05"))
}