package nr_types
import "encoding/json"

type NRMetric []NRMetricElement

func UnmarshalNRMetric(data []byte) (NRMetric, error) {
	var r NRMetric
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NRMetric) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NRMetricElement struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Value      float64    `json:"value"`
	Timestamp  int64      `json:"timestamp"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	DeviceID string `json:"device.id"`
}
