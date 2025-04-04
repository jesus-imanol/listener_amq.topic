package entities


type SensorData struct {
	ID        int64   `json:"id"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
	Text      string  `json:"text"`
	User      string  `json:"user"`
	CreatedAt string `json:"created_at"`
}
