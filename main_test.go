package main

import "testing"

func TestnotificationWorker(t *testing.T) {
	payload := Payload{
		Text:      "test",
		Username:  "TestBot",
		Channel:   "#test",
	}
	err := send("testUrl", "", payload)
	if err != nil {
		t.Fatal(err)
	}
}


//func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
//}
//
//func send(webhookURL string, proxy string, payload Payload) []error {
//}
//
//func notificationFunc(payload Payload, webhookURL string) []error {
//}
//
//func notificationWorker(queue string, args ...interface{}) error {
//}