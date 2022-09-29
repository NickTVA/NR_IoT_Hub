package nr_types

type PingEvent struct {
	EventType string `json:"eventType"`
	DeviceId  string `json:"device.id"`
}

type EnvironmentGeo struct {
	EventType   string  `json:"eventType"`
	DeviceId    string  `json:"device.id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
