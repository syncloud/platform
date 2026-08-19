package model

const (
	DiskSpaceSystem = "system"
	DiskSpaceData   = "data"
)

type DiskSpace struct {
	Low    bool             `json:"low"`
	Mounts []DiskSpaceMount `json:"mounts"`
}

type DiskSpaceMount struct {
	Kind    string `json:"kind"`
	Path    string `json:"path"`
	TotalKB uint64 `json:"total_kb"`
	FreeKB  uint64 `json:"free_kb"`
	Low     bool   `json:"low"`
}
