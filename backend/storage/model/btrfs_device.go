package model

type BtrfsDevice struct {
	Device string
	Devid  string
	Errors DeviceErrors
}
