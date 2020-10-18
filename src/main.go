package main

import (
	//"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack" //you must run 'go get github.com/nlopes/slack' in your terminal for this import to work
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"strconv"
	"time"
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
	starWarsMap := make(map[string]int)
	starWarsMap["yoda"] = 0
	starWarsMap["sith"] = 1
	starWarsMap["gungan"] = 2
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
						if strings.EqualFold(correctStr, "hello") { //simple hello command example
							rtm.SendMessage(rtm.NewOutgoingMessage("Hi", ev.Channel)) //sending Hi back in the channel
						} else if strings.EqualFold(correctStr, "joke") { //a little bit more complex joke command
							url := "https://official-joke-api.appspot.com/jokes/random" //url for api call (check it out, once you see it you'll understand why the struct below is important)
							resp, _ := http.Get(url) //http request to GET this URL
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
							rtm.SendMessage(rtm.NewOutgoingMessage(joke.Punchline, ev.Channel)) //send the second part of the joke
						} else if strings.EqualFold(correctStr, "dog"){
							url := "https://api.thedogapi.com/v1/images/search"
							resp, _ := http.Get(url)
							encodedDog, _ := ioutil.ReadAll(resp.Body)
							jsonDog := string(encodedDog)
							//https://mholt.github.io/json-to-go/
							type Weight struct {
								Imperial string `json:"imperial"`
								Metric   string `json:"metric"`
							}
							type Height struct {
								Imperial string `json:"imperial"`
								Metric   string `json:"metric"`
							}
							type Breeds struct {
								Weight      Weight `json:"weight"`
								Height      Height `json:"height"`
								ID          int    `json:"id"`
								Name        string `json:"name"`
								CountryCode string `json:"country_code"`
								BredFor     string `json:"bred_for"`
								BreedGroup  string `json:"breed_group"`
								LifeSpan    string `json:"life_span"`
								Temperament string `json:"temperament"`
							}
							type Dog []struct {
								Breeds []Breeds `json:"breeds"`
								ID     string   `json:"id"`
								URL    string   `json:"url"`
								Width  int      `json:"width"`
								Height int      `json:"height"`
							}
							var dog Dog
							json.Unmarshal([]byte(jsonDog), &dog)
							rtm.SendMessage(rtm.NewOutgoingMessage(dog[0].URL, ev.Channel))
							//rtm.SendMessage(rtm.NewOutgoingMessage(dog[0].Breeds[0].Name + " " + dog[0].Breeds[0].LifeSpan + " " + dog[0].Breeds[0].Weight.Metric, ev.Channel))
						} else if strings.EqualFold(correctStr, "kanye") {
							url := "https://api.kanye.rest"
							resp, _ := http.Get(url)

							encodedKanyeQuote, _ := ioutil.ReadAll(resp.Body)

							jsonKanyeQuote := string(encodedKanyeQuote)

							type KanyeQuote struct {
								Quote string `json:"quote"`
							}

							var kanyeQuote KanyeQuote
							json.Unmarshal([]byte(jsonKanyeQuote), &kanyeQuote)

							rtm.SendMessage(rtm.NewOutgoingMessage(kanyeQuote.Quote, ev.Channel))
						} else if strings.HasPrefix(strings.ToLower(correctStr), strings.ToLower("create channel")){
							channelNameAsArr := strings.Fields(correctStr)[2:]
							channelName := strings.Join(channelNameAsArr, " ")
							err, _ := rtm.CreateChannel(channelName)
							fmt.Println(err)
						} else if strings.HasPrefix(strings.ToLower(correctStr), strings.ToLower("rates")) {
							entireStr := strings.Split(correctStr, " ")
							baseCurrency := strings.ToUpper(entireStr[1])

							url := "https://api.exchangeratesapi.io/latest?base=" + baseCurrency
							resp, _ := http.Get(url)

							encodedRates, _ := ioutil.ReadAll(resp.Body)

							//fmt.Printf("%T\n", jsonRates)

							jsonRates := string(encodedRates)

							type RateResponse struct {
								Rates struct {
									CAD float64 `json:"CAD"`
									HKD float64 `json:"HKD"`
									ISK float64 `json:"ISK"`
									PHP float64 `json:"PHP"`
									DKK float64 `json:"DKK"`
									HUF float64 `json:"HUF"`
									CZK float64 `json:"CZK"`
									AUD float64 `json:"AUD"`
									RON float64 `json:"RON"`
									SEK float64 `json:"SEK"`
									IDR float64 `json:"IDR"`
									INR float64 `json:"INR"`
									BRL float64 `json:"BRL"`
									RUB float64 `json:"RUB"`
									HRK float64 `json:"HRK"`
									JPY float64 `json:"JPY"`
									THB float64 `json:"THB"`
									CHF float64 `json:"CHF"`
									SGD float64 `json:"SGD"`
									PLN float64 `json:"PLN"`
									BGN float64 `json:"BGN"`
									TRY float64 `json:"TRY"`
									CNY float64 `json:"CNY"`
									NOK float64 `json:"NOK"`
									NZD float64 `json:"NZD"`
									ZAR float64 `json:"ZAR"`
									USD float64 `json:"USD"`
									MXN float64 `json:"MXN"`
									ILS float64 `json:"ILS"`
									GBP float64 `json:"GBP"`
									KRW float64 `json:"KRW"`
									MYR float64 `json:"MYR"`
								} `json:"rates"`
								Base string `json:"base"`
								Date string `json:"date"`
							}


							var rateResponse RateResponse



							json.Unmarshal([]byte(jsonRates), &rateResponse)

							// Brute Force lmao
							rateMsg := "1.00 " + baseCurrency + " is \n" +
								strconv.FormatFloat(rateResponse.Rates.CAD, 'f', 6, 64) + " CAD\n" +
								strconv.FormatFloat(rateResponse.Rates.HKD, 'f', 6, 64) + " HKD\n" +
								strconv.FormatFloat(rateResponse.Rates.ISK, 'f', 6, 64) + " ISK\n" +
								strconv.FormatFloat(rateResponse.Rates.DKK, 'f', 6, 64) + " DKK\n" +
								strconv.FormatFloat(rateResponse.Rates.CZK, 'f', 6, 64) + " CZK\n" +
								strconv.FormatFloat(rateResponse.Rates.RON, 'f', 6, 64) + " RON\n" +
								strconv.FormatFloat(rateResponse.Rates.SEK, 'f', 6, 64) + " SEK\n" +
								strconv.FormatFloat(rateResponse.Rates.IDR, 'f', 6, 64) + " IDR\n" +
								strconv.FormatFloat(rateResponse.Rates.INR, 'f', 6, 64) + " INR\n" +
								strconv.FormatFloat(rateResponse.Rates.BRL, 'f', 6, 64) + " BRL\n"

							rtm.SendMessage(rtm.NewOutgoingMessage(rateMsg, ev.Channel))

					} else if _, exists := starWarsMap[strings.Fields(strings.ToLower(correctStr))[0]]; exists {
							loweredStr := strings.ToLower(correctStr)
							url := "https://api.funtranslations.com/translate/" + strings.Fields(loweredStr)[0] + ".json?text=" + strings.Join(strings.Fields(loweredStr)[1:], "%20")
							resp, _ := http.Get(url)
							encodedQuote, _ := ioutil.ReadAll(resp.Body)
							jsonQuotes := string(encodedQuote)
							type StarWarsQuote struct {
								Success struct {
									Total int `json:"total"`
								} `json:"success"`
								Contents struct {
									Translated  string `json:"translated"`
									Text        string `json:"text"`
									Translation string `json:"translation"`
								} `json:"contents"`
							}
							var StarWars StarWarsQuote
							json.Unmarshal([]byte(jsonQuotes), &StarWars)
							rtm.SendMessage(rtm.NewOutgoingMessage("Text: " + StarWars.Contents.Text, ev.Channel))
							rtm.SendMessage(rtm.NewOutgoingMessage("Translated: " + StarWars.Contents.Translated, ev.Channel))
					} else if strings.EqualFold(correctStr, "err") {
							codes := [10]int{200, 201, 204, 304, 400, 401, 403, 404, 409, 500}
							randNum := rand.Intn(len(codes))
							link := "https://http.cat/" + strconv.Itoa(codes[randNum]) + ".jpg"
							rtm.SendMessage(rtm.NewOutgoingMessage(link, ev.Channel))
					} else if strings.HasPrefix(strings.ToLower(correctStr), strings.ToLower("weather")){
							city := strings.Fields(correctStr)[1]
							preUrl := "https://www.metaweather.com/api/location/search/?query=" + city
							preResp, _ := http.Get(preUrl)
							encodedWOE, _ := ioutil.ReadAll(preResp.Body)
							jsonWOE := string(encodedWOE)
							type WOE []struct {
								Title        string `json:"title"`
								LocationType string `json:"location_type"`
								Woeid        int    `json:"woeid"`
								LattLong     string `json:"latt_long"`
							}
							var woe WOE
							json.Unmarshal([]byte(jsonWOE), &woe)
							woeID := woe[0].Woeid
							url := "https://www.metaweather.com/api/location/" + strconv.Itoa(woeID)
							resp, _ := http.Get(url)
							encodedData, _ := ioutil.ReadAll(resp.Body)
							jsonData := string(encodedData)
							type WeatherData struct {
								ConsolidatedWeather []struct {
									ID                   int64     `json:"id"`
									WeatherStateName     string    `json:"weather_state_name"`
									WeatherStateAbbr     string    `json:"weather_state_abbr"`
									WindDirectionCompass string    `json:"wind_direction_compass"`
									Created              time.Time `json:"created"`
									ApplicableDate       string    `json:"applicable_date"`
									MinTemp              float64   `json:"min_temp"`
									MaxTemp              float64   `json:"max_temp"`
									TheTemp              float64   `json:"the_temp"`
									WindSpeed            float64   `json:"wind_speed"`
									WindDirection        float64   `json:"wind_direction"`
									AirPressure          float64   `json:"air_pressure"`
									Humidity             int       `json:"humidity"`
									Visibility           float64   `json:"visibility"`
									Predictability       int       `json:"predictability"`
								} `json:"consolidated_weather"`
								Time         string `json:"time"`
								SunRise      string `json:"sun_rise"`
								SunSet       string `json:"sun_set"`
								TimezoneName string `json:"timezone_name"`
								Parent       struct {
									Title        string `json:"title"`
									LocationType string `json:"location_type"`
									Woeid        int    `json:"woeid"`
									LattLong     string `json:"latt_long"`
								} `json:"parent"`
								Sources []struct {
									Title     string `json:"title"`
									Slug      string `json:"slug"`
									URL       string `json:"url"`
									CrawlRate int    `json:"crawl_rate"`
								} `json:"sources"`
								Title        string `json:"title"`
								LocationType string `json:"location_type"`
								Woeid        int    `json:"woeid"`
								LattLong     string `json:"latt_long"`
								Timezone     string `json:"timezone"`
							}
							var weather WeatherData
							json.Unmarshal([]byte(jsonData), &weather)
							rtm.SendMessage(rtm.NewOutgoingMessage(strings.ToUpper(city) + " Weather: " + weather.ConsolidatedWeather[0].WeatherStateName + "\n" +
								"Temperature Now: " + fmt.Sprintf("%.2f", ((weather.ConsolidatedWeather[0].TheTemp) * 9/5) + 32) + " F\n" +
								"Max Temperature Today: " + fmt.Sprintf("%.2f", ((weather.ConsolidatedWeather[0].MaxTemp) * 9/5) + 32) + " F\n" +
								"Min Temperature Today: " + fmt.Sprintf("%.2f", ((weather.ConsolidatedWeather[0].MinTemp) * 9/5) + 32) + " F", ev.Channel))
						} else if strings.HasPrefix(strings.ToLower(correctStr), strings.ToLower("lyrics")) {
							entireStr := strings.Split(correctStr, ",")

							artist := strings.ReplaceAll(entireStr[1], " ", "%20")

							songName := strings.ReplaceAll(entireStr[2], " ", "%20")

							url := "https://api.lyrics.ovh/v1/" + artist + "/" + songName
							resp, _ := http.Get(url)

							encodedLyrics, _ := ioutil.ReadAll(resp.Body)
							jsonLyrics := string(encodedLyrics)

							type SongLyrics struct {
								Lyrics string `json:"lyrics"`
							}

							var lyrics SongLyrics
							json.Unmarshal([]byte(jsonLyrics), &lyrics)

							rtm.SendMessage(rtm.NewOutgoingMessage(lyrics.Lyrics, ev.Channel))
						} else if strings.EqualFold(correctStr, "help"){
							rtm.SendMessage(rtm.NewOutgoingMessage("!hello - says Hi\n!joke - returns joke\n!dog - return dog pic" +
								"\n!kanye - returns kanye quote\n!rates - converts 1 currency to 10 different currencies" +
								"\n![yoda, sith, gungan] <text> - returns text translated into any of the 3 characters from Star Wars" +
								"\n!err - returns a random cat pic from 10 error codes\n!weather <city> - returns brief details about weather in entered city" +
								"\n!lyrics, <artist name>, <artist song> - returns lyrics by an artist" +
								"\n!help - returns this command", ev.Channel))
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
