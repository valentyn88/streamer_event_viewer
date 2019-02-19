package streamer_event_viewer

// SubscriptionBody subscription body.
type SubscriptionBody struct {
	Callback string `json:"hub.callback"`
	Mode     string `json:"hub.mode"`
	Topic    string `json:"hub.topic"`
	Seconds  int    `json:"hub.lease_seconds"`
}
