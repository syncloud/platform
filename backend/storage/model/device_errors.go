package model

type DeviceErrors struct {
	Write      int64 `json:"write"`
	Read       int64 `json:"read"`
	Flush      int64 `json:"flush"`
	Corruption int64 `json:"corruption"`
	Generation int64 `json:"generation"`
}

func (e DeviceErrors) Any() bool {
	return e.Write != 0 || e.Read != 0 || e.Flush != 0 || e.Corruption != 0 || e.Generation != 0
}
