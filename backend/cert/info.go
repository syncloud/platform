package cert

type Info struct {
	Subject      string `json:"-"`
	IsValid      bool   `json:"is_valid"`
	IsReal       bool   `json:"is_real"`
	ValidForDays int    `json:"valid_for_days"`
}
