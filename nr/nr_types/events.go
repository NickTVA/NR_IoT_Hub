package nr_types

type PingEvent struct {
	EventType string `json:"eventType"`
	DeviceId  string `json:"device.id"`
}
