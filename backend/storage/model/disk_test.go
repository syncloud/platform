package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindRootPartitionSome(t *testing.T) {
	disk := Disk{"disk", "/dev/sda", "20", []Partition{
		{"10", "/dev/sda1", "/", true, "ext4", false},
		{"10", "/dev/sda2", "", true, "ext4", false},
	}, true, "", "", "", false, false}

	assert.Equal(t, disk.FindRootPartition().Device, "/dev/sda1")
}

func TestFindRootPartition_Nil(t *testing.T) {
	disk := Disk{"disk", "/dev/sda", "20", []Partition{
		{"10", "/dev/sda1", "/my", true, "ext4", false},
		{"10", "/dev/sda2", "", true, "ext4", false},
	}, true, "", "", "", false, false}
	assert.Nil(t, disk.FindRootPartition())
}
