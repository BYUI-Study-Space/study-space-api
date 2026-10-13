// Package models holds the API's data types and the session validation rules.
package models

// Location is a study spot. Field order matches the store's SELECT column order.
type Location struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Building      string `json:"building"`
	Room          string `json:"room"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	HasWifi       bool   `json:"has_wifi"`
	GroupFriendly bool   `json:"group_friendly"`
	NoiseLevel    string `json:"noise_level"`
	Capacity      int    `json:"capacity"`
	OpenTime      string `json:"open_time"`  // "HH:MM"
	CloseTime     string `json:"close_time"` // "HH:MM"
}

// LocationFilter narrows the location list; nil or empty fields don't filter.
type LocationFilter struct {
	Wifi  *bool
	Group *bool
	Noise string
}

var NoiseLevels = map[string]bool{"quiet": true, "moderate": true, "collaborative": true}
