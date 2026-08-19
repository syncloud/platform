package model

type SnapshotSet struct {
	Id        int        `json:"id"`
	Snapshots []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Set  int    `json:"set"`
	Snap string `json:"snap"`
	Size uint64 `json:"size"`
	Auto bool   `json:"auto"`
}

func (s *SnapshotSet) IsAuto() bool {
	if len(s.Snapshots) == 0 {
		return false
	}
	for _, snapshot := range s.Snapshots {
		if !snapshot.Auto {
			return false
		}
	}
	return true
}

type SnapshotSetsResponse struct {
	Result []SnapshotSet `json:"result"`
}

type SnapshotsRequest struct {
	Action string `json:"action"`
	Set    int    `json:"set"`
}
