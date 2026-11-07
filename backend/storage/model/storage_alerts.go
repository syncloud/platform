package model

type StorageAlerts struct {
	Problems []StorageProblem `json:"problems"`
	New      bool             `json:"new"`
}
