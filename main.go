package main

import (
	"NR_IoT_Hub/nr/nr_types"
	"bytes"
	"fmt"
	insights "github.com/newrelic/go-insights/client"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const MetricUrl = "https://metric-api.newrelic.com/metric/v1"
const LogsUrl = "https://log-api.newrelic.com/log/v1"

func main() {

	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) < 2 {
		log.Fatal("usage: nr_insights_key nr_account_id [listenAddress]")
	}

	port := "4590"

	if len(cmdLineArgs) == 3 {
		port = cmdLineArgs[2]
	}

	apiKey := cmdLineArgs[0]

	account_id := cmdLineArgs[1]

	log.Println("Starting NR IOT Hub on port ..." + port)
	log.Println("Insights key: " + apiKey)
	log.Println("Account id: " + account_id)

	http.HandleFunc("/metric", handleMetric(apiKey))
	http.HandleFunc("/log", handleLog(apiKey))
	http.HandleFunc("/ping", handlePing(apiKey, account_id))

	listenAddress := ":" + port
	log.Fatal(http.ListenAndServe(listenAddress, nil))

}

func handlePing(apiKey string, account_id string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		id_array, exists := query["id"]
		if !exists || (len(id_array[0]) < 1) {
			_, _ = fmt.Fprintf(w, "No device id_array")
			return

		}

		id := id_array[0]

		type_array, exists := query["type"]

		if !exists || (len(type_array[0]) < 1) {

			_, _ = fmt.Fprintf(w, "No  type")
			return
		}

		eventType := type_array[0]

		event := nr_types.PingEvent{
			EventType: eventType,
			DeviceId:  id,
		}

		sendEvent(event, apiKey, account_id)

	}
}

func handleMetric(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		id_array, exists := query["id"]
		if !exists || (len(id_array[0]) < 1) {
			_, _ = fmt.Fprintf(w, "No device id_array")
			return

		}

		id := id_array[0]

		name_array, exists := query["name"]
		if !exists || (len(name_array[0]) < 1) {
			_, _ = fmt.Fprintf(w, "No name")
			return
		}

		name := name_array[0]

		type_array, exists := query["type"]

		if !exists || (len(type_array[0]) < 1) {

			_, _ = fmt.Fprintf(w, "No  type")
			return
		}

		metricType := type_array[0]

		value_array, exists := query["value"]

		if !exists || (len(value_array[0]) < 1) {

			_, _ = fmt.Fprintf(w, "No value_array")
			return
		}

		value := value_array[0]
		value_float, _ := strconv.ParseFloat(value, 32)

		nrmetric := makeMetric(id, name, metricType, value_float)

		response, err := sendNRMetric(nrmetric, apiKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println(response)

		log.Println(response.Status)

		w.Write([]byte("NR status: " + string(response.Status)))

	}
}

func handleLog(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		id_array, exists := query["id"]
		if !exists || (len(id_array[0]) < 1) {
			_, _ = fmt.Fprintf(w, "No device id_array")
			return
		}

		id := id_array[0]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		body_string := string(body)

		nrLog := makeLog(id, body_string)

		response, err := sendNRLog(nrLog, apiKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println(response)

		log.Println(response.Status)
		w.Write([]byte("NR log status: " + string(response.Status)))

	}
}

func sendNRMetric(nrmetric nr_types.NRMetric, apiKey string) (*http.Response, error) {
	client := http.Client{}
	bts, _ := nrmetric.Marshal()
	req, _ := http.NewRequest("POST", MetricUrl, bytes.NewBuffer(bts))
	req.Header.Add("Api-Key", apiKey)
	log.Println("Sending NR Metric...")
	log.Println(string(bts))
	response, err := client.Do(req)
	return response, err
}

func makeMetric(id string, name string, metricType string, value_float float64) nr_types.NRMetric {
	metric_attributes := nr_types.Attributes{DeviceID: id}

	curr_time_millis := time.Now().Unix()

	metric := nr_types.Metric{
		Name:       name,
		Type:       metricType,
		Value:      value_float,
		Timestamp:  curr_time_millis,
		Attributes: metric_attributes,
	}

	nrmetricsElements := nr_types.NRMetricElement{Metrics: []nr_types.Metric{metric}}

	nrmetric := nr_types.NRMetric{nrmetricsElements}
	return nrmetric
}

func sendNRLog(nrlog nr_types.NRLog, apiKey string) (*http.Response, error) {
	client := http.Client{}
	bts, _ := nrlog.Marshal()
	req, _ := http.NewRequest("POST", LogsUrl, bytes.NewBuffer(bts))
	req.Header.Add("X-Insert-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")
	log.Println("Sending NR Log...")
	log.Println(string(bts))
	response, err := client.Do(req)
	return response, err
}

func sendEvent(event interface{}, apikey string, insightAccountID string) {
	client := insights.NewInsertClient(apikey, insightAccountID)
	if validationErr := client.Validate(); validationErr != nil {
		//however it is appropriate to handle this in your use case
		log.Println("Insights Client Validation Error!")
	}

	if postErr := client.PostEvent(event); postErr != nil {
		log.Println("Error: %v\n", postErr)
	}
	log.Println(event)
	log.Println(client)
}

func makeLog(id string, message string) nr_types.NRLog {

	curr_time_millis := time.Now().Unix()
	logType := "IoT"

	log := nr_types.NRLog{
		Timestamp: curr_time_millis,
		Message:   message,
		Logtype:   logType,
		DeviceID:  id,
	}

	return log

}
