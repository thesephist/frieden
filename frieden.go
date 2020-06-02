package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// static file server
// api endpoint that runs as a (cached?) proxy to gcal servers

var config appConfig

type appConfig struct {
	ApiKey      string   `json:"apiKey"`
	CalendarIds []string `json:"calendars"`
}

type calendarIncomingReq struct {
	TimeZone string `json:"timeZone"`
	TimeMin  string `json:"timeMin"`
	TimeMax  string `json:"timeMax"`
}

type calendarOutgoingReq struct {
	TimeZone  string                  `json:"timeZone"`
	TimeMin   string                  `json:"timeMin"`
	TimeMax   string                  `json:"timeMax"`
	Calendars []calendarOutgoingReqId `json:"items"`
}

type calendarOutgoingReqId struct {
	Id string `json:"id"`
}

func mustConfigure() {
	configFile, err := os.Open("./secrets.json")
	if err != nil {
		log.Fatalf("Could not read secrets file!\n\tError: %s", err.Error())
	}
	defer configFile.Close()

	configString, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatalf("Could not read config from file!\n\tError: %s", err.Error())
	}
	err = json.Unmarshal(configString, &config)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	indexFile, err := os.Open("./static/index.html")
	if err != nil {
		io.WriteString(w, "error reading index")
		return
	}
	defer indexFile.Close()

	io.Copy(w, indexFile)
}

func getData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, "{}")
		return
	}

	var reqParams calendarIncomingReq
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		io.WriteString(w, "{}")
		return
	}

	calendars := []calendarOutgoingReqId{}
	for _, id := range config.CalendarIds {
		calendars = append(calendars, calendarOutgoingReqId{
			Id: id,
		})
	}
	outgoingParams := calendarOutgoingReq{
		TimeZone:  reqParams.TimeZone,
		TimeMin:   reqParams.TimeMin,
		TimeMax:   reqParams.TimeMax,
		Calendars: calendars,
	}

	jsonStr, err := json.Marshal(&outgoingParams)
	proxyReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://www.googleapis.com/calendar/v3/freeBusy?alt=json&key=%s", config.ApiKey),
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		// TODO
		log.Fatal("TODO")
	}
	proxyReq.Header.Set("Accept", "application/json")
	proxyReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		// TODO
		log.Fatal("TODO")
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// TODO
		log.Fatal("TODO")
	}
}

func main() {
	mustConfigure()

	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:7856",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	r.HandleFunc("/", handleHome)
	r.Methods("POST").Path("/data").HandlerFunc(getData)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Frieden listening on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
