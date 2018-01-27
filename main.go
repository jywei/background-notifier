package main

import (
	// "bufio"
	"fmt"
	// "os"

	"github.com/garyburd/redigo/redis"
	"github.com/parnurzeal/gorequest"
	"github.com/gin-gonic/gin"
)

// Payload Type is for saving payload as JSON strcture
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

func notificationWorker(webhookURL string, response chan<- []error) {
	payload := Payload{
		Text:      ":thinking_face:",
		Username:  "Hexter Bot",
		IconEmoji: ":hexter_is_ur_daddy:",
		Channel:   "#studygroup-tw",
	}
	response <- send(webhookURL, "", payload)
}

func main() {
	c, err := redis.Dial("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	defer c.Close()


	// There code is for console testing


	//for {
	//	fmt.Println(">>Please enter some messages for notification Enter it or \"q\" to exit")
	//	bio := bufio.NewReader(os.Stdin)
	//	line, _, _ := bio.ReadLine()
	//
	//	if string(line) == "q" {
	//		break
	//	}
	//
	//	c.Do("SET", "notification", line)
	//	message, err := redis.String(c.Do("GET", "notification"))
	//	if err != nil {
	//		fmt.Println("key not found")
	//	}
	//	payload := Payload{
	//		Text:      message,
	//		Username:  "Hexter Bot",
	//		IconEmoji: ":hexter_is_ur_daddy:",
	//		Channel:   "#studygroup-tw",
	//	}
	//	sendErr := send(webhookURL, "", payload)
	//	if sendErr != nil {
	//		fmt.Printf("error: %s\n", err)
	//	}
	//}

	r := gin.Default()
	r.GET("/notification", func(c *gin.Context) {
		webhookURL := "https://hooks.slack.com/services/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"
		responseChannel := make(chan []error)
		go notificationWorker(webhookURL, responseChannel)
		err := <- responseChannel
		if err != nil {
			c.JSON(400, gin.H{
				"message": err,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Goodjob",
			})
		}
		close(responseChannel)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
