package streamer_event_viewer

// Data response from Twitch API.
type Data struct {
	Users []User `json:"data"`
}

// User twitch user.
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

// Streamer twitch streamer.
type Streamer struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}
