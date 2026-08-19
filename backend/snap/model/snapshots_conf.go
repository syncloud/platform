package model

import "encoding/json"

const RetentionKey = "snapshots.automatic.retention"

type SnapshotsConf struct {
	Snapshots SnapshotsConfSnapshots `json:"snapshots"`
}

type SnapshotsConfSnapshots struct {
	Automatic SnapshotsConfAutomatic `json:"automatic"`
}

type SnapshotsConfAutomatic struct {
	Retention string `json:"retention"`
}

func NewSnapshotsConf(retention string) SnapshotsConf {
	return SnapshotsConf{
		Snapshots: SnapshotsConfSnapshots{
			Automatic: SnapshotsConfAutomatic{Retention: retention},
		},
	}
}

type ConfResponse struct {
	Type   string          `json:"type"`
	Result json.RawMessage `json:"result"`
}

func (r *ConfResponse) Value(key string) string {
	if r.Type != "sync" {
		return ""
	}
	var values map[string]string
	err := json.Unmarshal(r.Result, &values)
	if err != nil {
		return ""
	}
	return values[key]
}
