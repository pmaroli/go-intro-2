package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"time"
)

// PORT is the port on which the server runs
const PORT = ":12345"

// Global variables
var euler float64 = math.E
var totalRequests uint32 = 0
var startTime string
var lastFetchedTime string

// Make the channel that logs to the 'logs' file
var logger = make(chan loggerData)

// Create structs
type worldTimeAPI struct {
	DateTime string `json:"datetime"`
}

type loggerData struct {
	RequestIP       string
	LastFetchedTime string
	RequestTime     string
}

func main() {
	startTime = time.Now().Format(time.RFC3339)
	lastFetchedTime = startTime
	// ticker creates a channel ticker.C which is sent a tick every 2.7~ seconds
	var ticker = time.NewTicker(time.Duration(int(euler * 1000000000)))

	// Start the background GoRoutine that logs to the 'logs' file
	go Logger(logger)

	// Start the background GoRoutine that fetches the time every e seconds
	go GetTime(ticker)

	http.HandleFunc("/", Root)

	fmt.Printf("Server is starting on PORT %s ...\n", PORT[1:])
	http.ListenAndServe(PORT, nil)
}

// Root handles all requests to the root endpoint
func Root(w http.ResponseWriter, req *http.Request) {
	totalRequests = totalRequests + 1
	fmt.Fprintf(w, "Start Time: %s\n", startTime)
	fmt.Fprintf(w, "Last Fetched Time: %s\n", lastFetchedTime)
	fmt.Fprintf(w, "Number of Requests Made: %d\n", totalRequests)

	// Parse the client's IP address
	host, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Println("Could not get IP Address!\n", err)
	}

	// Send the new loggerData to the logger channel
	logger <- loggerData{
		RequestIP:       host + ":" + port,
		LastFetchedTime: lastFetchedTime,
		RequestTime:     time.Now().Format(time.RFC3339),
	}
}

// Logger logs to the 'logs' file every time new loggerData is added to the logger channel
func Logger(logger chan loggerData) {
	for {
		var log = <-logger

		// Append to 'logs' file
		var newLog = fmt.Sprintf("<%s>-<%s>-<%s>\n", log.RequestIP, log.LastFetchedTime, log.RequestTime)

		f, err := os.OpenFile("logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error while opening file!\n", err)
		}
		defer f.Close()
		if _, err := f.WriteString(newLog); err != nil {
			fmt.Println("Error while writing to file!\n", err)
		}
	}
}

// GetTime gets time data from worldtimeapi.org every E seconds
func GetTime(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			resp, err := http.Get("http://worldtimeapi.org/api/ip")
			if err != nil {
				fmt.Println("Error fetching the time from worldtimeapi.org!\n", err)
			}
			defer resp.Body.Close()

			respBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading the response body!\n", err)
			}

			var respJSON worldTimeAPI

			err = json.Unmarshal(respBytes, &respJSON)
			if err != nil {
				fmt.Println("There was an error parsing the JSON!\n", err)
			}

			tempTime, err := time.Parse(time.RFC3339, respJSON.DateTime)
			if err != nil {
				fmt.Println("An error occurred while parsing the time!\n", err)
			}

			// Parse the tempTime time.Time variable into a string
			lastFetchedTime = tempTime.Format(time.RFC3339)
		}
	}
}
