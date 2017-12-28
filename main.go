package main

import (
	"bytes"
	"fmt"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"log"
	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	configFile string = "config.yml"
)

type BotMessage struct {
	Channel string `json:"channel,omitempty"`
	User    string `json:"user,omitempty"`
	Text    string `json:"text,omitempty"`
	Team    string `json:"team,omitempty"`
}

type Config struct {
	Api Api
}

type Api struct {
	Key string `key`
}

var (
	botId string
)

func isMentioned(event *slack.MessageEvent) bool {
	if event.Type == "message" && strings.HasPrefix(event.Text, "<@"+botId+">") {
		return true
	}
	return false
}

func radioMessage(air OnAir) string {
	var buffer bytes.Buffer

	buffer.WriteString("Nu op Top2000: " + air.Results[0].Songfile.Title + " -  " + air.Results[0].Songfile.Artist)

	return buffer.String()
}

func handleRadioReport(command []string, message BotMessage, reply chan<- BotMessage) {
	if len(command) > 1 {
		data, err := getAapjeData()
		if err != nil {
			log.Fatal(err)
		}

		message.Text = radioMessage(data)

		fmt.Printf("Sending reply to channel\n")
		reply <- message
		fmt.Printf("Reply sent to channel\n")
	}
}

func BotMentionChannelHandler(incoming <-chan BotMessage, reply chan<- BotMessage) {
	fmt.Printf("Started bot handler channel\n")
	for message := range incoming {
		fmt.Printf("Handling mention\n")

		command := strings.Fields(message.Text)

		if len(command) > 0 {
			myCommand := command[1]

			switch {
			case strings.Contains(myCommand, "nu"):
				go handleRadioReport(command, message, reply)
			}
		}
	}
}

func BotReplyChannelHandler(incoming <-chan BotMessage, api *slack.Client) {
	fmt.Printf("Starting bot reply handler channel\n")

	for message := range incoming {
		fmt.Printf("Handling reply\n")
		//spew.Dump(message)

		replyParameters := slack.PostMessageParameters{}
		replyParameters.AsUser = true

		fmt.Printf("Posting reply")
		_, _, err := api.PostMessage(message.Channel, message.Text, replyParameters)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readConfig() Config {
	fmt.Printf("Obtaining config\n")

	var config Config
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reading configfile\n")
	yaml.Unmarshal(source, &config)

	return config
}

func getToken() string {
	config := readConfig()
	fmt.Printf("Returning API token\n")
	return config.Api.Key
}

func main() {
	api := slack.New(getToken())
	rtm := api.NewRTM()

	botMentionChannel := make(chan BotMessage)
	botReplyChannel := make(chan BotMessage)

	go rtm.ManageConnection()
	go BotMentionChannelHandler(botMentionChannel, botReplyChannel)
	go BotReplyChannelHandler(botReplyChannel, api)

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch event := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connected!\n")
				botId = event.Info.User.ID

			case *slack.MessageEvent:
				if isMentioned(event) {
					fmt.Printf("Was mentioned!\n")

					fmt.Printf("Sending to channel!\n")
					botMentionChannel <- BotMessage{
						Channel: event.Channel,
						User:    event.User,
						Text:    event.Text,
						Team:    event.Team,
					}
				}
			}
		}
	}
}
