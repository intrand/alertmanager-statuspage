package main

import (
	"encoding/json" // this is a web service, so obviously we'll need JSON :)
	"io/ioutil"
	"log"
	"net/http" // server & client (both)
	"os"

	"gopkg.in/alecthomas/kingpin.v2" // very easy flag/env/defaults management
)

// a single alert
type alertManAlert struct {
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

// the entire alertmanager payload (including the list of alerts)
type alertManOut struct {
	Alerts            []alertManAlert `json:"alerts"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
	} `json:"commonLabels"`
	ExternalURL string `json:"externalURL"`
	GroupKey    string `json:"groupKey"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

// the input to the statuspage.io API
type statusPageIn struct {
	Component struct {
		Status string `json:"status"`
	} `json:"component"`
}

// constants
const apiUrl = "https://api.statuspage.io/v1" // base URL of the statuspage.io API

// flags, environment variables and defaults
var (
	ams           = kingpin.New("alertmanager-statuspage", "Takes alerts from alertmanager and updates statuspage.io accordingly.")
	token         = ams.Flag("token", "statuspage.io API token").Envar("token").Required().String()
	listenAddress = ams.Flag("listen.address", "address:port to listen on").Envar("listen_address").Default("127.0.0.1:8080").String()
)

// logic to determine if we need to send a request to the statuspage.io API or not
func filterAlerts(amo *alertManOut, token string) {
	for _, alert := range amo.Alerts {
		// guarantee we have a statuspageio_component label
		if _, ok := alert.Labels["statuspageio_component"]; !ok { // statuspageio_component label not configured for this alert
			continue // do nothing, then move on to the next one
		}

		// guarantee we have a statuspageio_page label
		if _, ok := alert.Labels["statuspageio_page"]; !ok { // statuspageio_page label not configured for this alert
			continue // do nothing, then move on to the next one
		}

		// guarantee we have a severity (default or found)
		severity := "under_maintenance"                             // assumed to be under maintenance and an admin forgot to set a silence :)
		if value, ok := alert.Labels["statuspageio_severity"]; ok { // but if we find a label indicating this is something else...
			severity = value // set the statuspage.io status for this component to what we found
		}

		println(alert.Annotations.Summary) // useful for debugging

		// build statuspage.io API request
		SPI := statusPageIn{}
		url := apiUrl + "/pages/" + alert.Labels["statuspageio_page"] + "/components/" + alert.Labels["statuspageio_component"]

		// status
		if alert.Status == "firing" {
			SPI.Component.Status = severity
		} else if alert.Status == "resolved" {
			SPI.Component.Status = "operational"
		}

		// Golang -> JSON
		SPID, _ := json.Marshal(SPI)

		// send HTTP PATCH to update the component on the page
		patchPage(SPID, url, token) // see http.go
	}
}

func main() {
	kingpin.MustParse(ams.Parse(os.Args[1:])) // get config
	*token = "OAuth " + *token                // format token appropriately

	log.Printf("Listening on: %s", *listenAddress)                                                      // useful for debugging
	http.ListenAndServe(*listenAddress, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // handle webhooks from alertmanager
		log.Printf("%s - [%s] %s", r.Host, r.Method, r.URL.RawPath) // useful for debugging

		b, err := ioutil.ReadAll(r.Body) // read the request
		if err != nil {                  // check for errors
			panic(err) // panic on error
		}

		// JSON -> Golang
		amo := alertManOut{}
		err = json.Unmarshal(b, &amo)
		if err != nil { // check for errors
			if len(b) > 1024 {
				log.Printf("Failed to unpack inbound alert request - %s...", string(b[:1023]))
			} else {
				log.Printf("Failed to unpack inbound alert request - %s", string(b))
			}

			return
		}

		filterAlerts(&amo, *token) // process valid alerts for potential statuspage.io changes
	}))
}
