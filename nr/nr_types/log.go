package nr_types

import "encoding/json"

func UnmarshalNRLog(data []byte) (NRLog, error) {
	var r NRLog
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NRLog) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NRLog struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Logtype   string `json:"logtype"`
	DeviceID  string `json:"device.id"`
}
