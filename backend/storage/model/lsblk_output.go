package model

type LsblkOutput struct {
	BlockDevices []LsblkDevice `json:"blockdevices"`
}
