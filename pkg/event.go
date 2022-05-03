package sentry

type Event struct {
	Timestamp string `json:"timestamp"`
	City      string `json:"city"`
	EstLat    string `json:"est_lat"`
	EstLon    string `json:"est_lon"`
	Reason    string `json:"reason"`
	Camera    string `json:"camera"`
}
