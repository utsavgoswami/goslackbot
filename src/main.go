package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	prefix := "!"
	contentAsArray, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatal(err)
	}
	botToken := string(contentAsArray)
	api := slack.New(botToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				command := ev.Text
				if strings.HasPrefix(command, prefix) {
					correctStr := command[1:]
					//todo: make this neater so that we call a method for each respective command instead of putting everything in here
					if ev.User != info.User.ID {
						if strings.EqualFold(correctStr, "hello") {
							rtm.SendMessage(rtm.NewOutgoingMessage("Hi", ev.Channel))
						} else if strings.EqualFold(correctStr, "joke"){

						}
					}
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Take no action
			}
		}
	}
}
