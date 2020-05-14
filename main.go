package main

import (
	"NR_IoT_Hub/nr/nr_types"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const MetricUrl = "https://metric-api.newrelic.com/metric/v1"
const LogsUrl = "https://log-api.newrelic.com/log/v1"

type NR_Metric []struct {
	Metrics []struct {
		Name      string  `json:"name"`
		Type      string  `json:"type"`
		Value     float64 `json:"value"`
		Timestamp int64   `json:"timestamp"`
	} `json:"metrics"`
}

func main() {

	cmdLineArgs := os.Args[1:]

	if len(cmdLineArgs) < 1 {
		log.Fatal("usage: nr_insights_key [listenAddress]")
	}

	port := "4590"

	if len(cmdLineArgs) == 2 {
		port = cmdLineArgs[1]
	}

	apiKey := cmdLineArgs[0]

	log.Println("Starting NR IOT Hub on port ..." + port)

	http.HandleFunc("/metric", handleMetric(apiKey))
	http.HandleFunc("/log", handleLog(apiKey))

	listenAddress := ":" + port
	log.Fatal(http.ListenAndServe(listenAddress, nil))

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
