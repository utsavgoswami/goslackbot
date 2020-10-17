package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack" //you must run 'go get github.com/nlopes/slack' in your terminal for this import to work
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("Bot is now running...")
	prefix := "!" //exclamation mark is prefix for all of our commands
	contentAsArray, err := ioutil.ReadFile("token.txt") //have a token.txt in your local dir with the token
	if err != nil {
		log.Fatal(err)
	}
	botToken := string(contentAsArray)

	//auth stuff
	api := slack.New(botToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	//

Loop: //begin
	for { //this for loop is necessary for the bot!
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent: //all of our commands will reside under this case
				info := rtm.GetInfo()
				command := ev.Text //user message
				if strings.HasPrefix(command, prefix) { //all commands MUST start with '!'. If a message doesn't start with '!' ignore it
					correctStr := command[1:]
					if ev.User != info.User.ID { //NECESSARY TO CHECK IF USER IS NOT THE BOT, otherwise recursion occurs
						//todo: try to make each command be its separate function? Idk, couldn't figure it out lol
						if strings.EqualFold(correctStr, "hello") { //simple hello command example
							rtm.SendMessage(rtm.NewOutgoingMessage("Hi", ev.Channel)) //sending Hi back in the channel
						} else if strings.EqualFold(correctStr, "joke") { //a little bit more complex joke command
							url := "https://official-joke-api.appspot.com/jokes/random" //url for api call (check it out, once you see it you'll understand why the struct below is important)
							resp, err := http.Get(url) //http request to GET this URL
							if err != nil {
								fmt.Println(err)
							}
							encodedJoke, _ := ioutil.ReadAll(resp.Body) //here, we read the response body (it is encoded)
							jsonJoke := string(encodedJoke) //once we preform the string function, joke turns into a json format
							type Joke struct { //this is necessary to decode the json, (https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/) good article
								ID        json.Number
								Type      string
								Setup     string
								Punchline string
							}
							var joke Joke //refer to previous comment
							json.Unmarshal([]byte(jsonJoke), &joke) //here we are decoding (unmarshal) the json and variable 'joke' from previous line will hold the result
							rtm.SendMessage(rtm.NewOutgoingMessage(joke.Setup, ev.Channel)) //send the first part of the joke
							//todo: some way to pause between messages?
							rtm.SendMessage(rtm.NewOutgoingMessage(joke.Punchline, ev.Channel)) //send the second part of the joke
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
} //end
