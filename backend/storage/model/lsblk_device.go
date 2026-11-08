package model

type LsblkDevice struct {
	Name       string        `json:"name"`
	Size       string        `json:"size"`
	Type       string        `json:"type"`
	MountPoint string        `json:"mountpoint"`
	PartType   string        `json:"parttype"`
	FsType     string        `json:"fstype"`
	Model      string        `json:"model"`
	Uuid       string        `json:"uuid"`
	Serial     string        `json:"serial"`
	Wwn        string        `json:"wwn"`
	Hctl       string        `json:"hctl"`
	Children   []LsblkDevice `json:"children"`
}
