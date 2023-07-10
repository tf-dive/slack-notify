package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	EnvSlackWebhook  			= "SLACK_WEBHOOK"
	EnvSlackIcon					= "SLACK_ICON"
	EnvSlackChannel  			= "SLACK_CHANNEL"
	EnvSlackTitle    			= "SLACK_TITLE"
	EnvSlackMessage  			= "SLACK_MESSAGE"
	EnvSlackAttachments  	= "SLACK_ATTACHMENTS"
	EnvSlackColor    			= "SLACK_COLOR"
	EnvSlackUserName 			= "SLACK_USERNAME"
)

type Webhook struct {
	Text        string       `json:"text,omitempty"`
	UserName    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback string  `json:"fallback"`
	Pretext  string  `json:"pretext"`
	Color    string  `json:"color"`
	Fields   []Field `json:"fields,omitempty"`
	MrkDwnIn []string `json:"mrkdwn_in"`
	Title    string  `json:"title"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value,omitempty"`
	Short bool   `json:short"`
}

func main() {
	endpoint := os.Getenv(EnvSlackWebhook)
	if endpoint == "" {
		fmt.Fprintln(os.Stderr, "URL is required")
		os.Exit(1)
	}

	test, found := os.LookupEnv("SLACK_ATTACHMENTS")
  if found {
      fmt.Println(test)
  } else {
      fmt.Println("$SLACK_ATTACHMENTS not found")
  }

	var attachments []Attachment
	attachmentsString := os.Getenv(EnvSlackAttachments)
	fmt.Fprintf(os.Stdout, "attachmentsString : %s\n", attachmentsString)
	if attachmentsString == "" {
		text := os.Getenv(EnvSlackMessage)
		fmt.Fprintf(os.Stdout, "text : %s\n", text)
		if text == "" {
			fmt.Fprintln(os.Stderr, "Message is required")
			os.Exit(1)
		}

		attachments = []Attachment{
			{
				Fallback: envOr(EnvSlackMessage, "This space intentionally left blank"),
				Color:    os.Getenv(EnvSlackColor),
				Fields: []Field{
					{
						Title: os.Getenv(EnvSlackTitle),
						Value: envOr(EnvSlackMessage, "EOM"),
					},
				},
			},
		}
	} else {
		attachments = []Attachment{}
		err := json.Unmarshal([]byte(attachmentsString), &attachments)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Attachments is wrong")
			fmt.Fprintln(os.Stderr, "err : ", err)
			os.Exit(1)
		}
	}

	if attachments == nil {
		fmt.Fprintln(os.Stderr, "Attachments is missing")
		os.Exit(1)
	}

	msg := Webhook{
		UserName: os.Getenv(EnvSlackUserName),
		IconURL:  os.Getenv(EnvSlackIcon),
		Channel:  os.Getenv(EnvSlackChannel),
		Attachments: attachments,
	}

	if err := send(endpoint, msg); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %s\n", err)
		os.Exit(2)
	}
}

func envOr(name, def string) string {
	if d, ok := os.LookupEnv(name); ok {
		return d
	}
	return def
}

func send(endpoint string, msg Webhook) error {
	enc, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(enc)
	res, err := http.Post(endpoint, "application/json", b)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		return fmt.Errorf("Error on message: %s\n", res.Status)
	}
	fmt.Println(res.Status)
	return nil
}
