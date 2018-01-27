package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/parnurzeal/gorequest"
)

// Payload Type is for ...
type Payload struct {
	Parse     string `json:"parse,omitempty"`
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

func main() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	webhookURL := "https://hooks.slack.com/services/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"

	for {
		fmt.Println(">>Please enter some messages for notification Enter it or \"q\" to exit")
		bio := bufio.NewReader(os.Stdin)
		line, _, _ := bio.ReadLine()

		if string(line) == "q" {
			break
		}

		c.Do("SET", "notification", line)
		message, err := redis.String(c.Do("GET", "notification"))
		if err != nil {
			fmt.Println("key not found")
		}
		payload := Payload{
			Text:      message,
			Username:  "Hexter Bot",
			IconEmoji: ":hexter_is_ur_daddy:",
		}
		sendErr := send(webhookURL, "", payload)
		if sendErr != nil {
			fmt.Printf("error: %s\n", err)
		}
	}
}
