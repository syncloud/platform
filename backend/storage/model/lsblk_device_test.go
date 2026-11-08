package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLsblkDevice_IsMMCBootPartition(t *testing.T) {
	assert.True(t, (&LsblkDevice{Name: "/dev/mmcblk1boot0"}).IsMMCBootPartition())
}

func TestLsblkDevice_IsMMCBootPartition_Not(t *testing.T) {
	assert.False(t, (&LsblkDevice{Name: "/dev/mmcblk1p0"}).IsMMCBootPartition())
}

func TestLsblkDevice_IsZram(t *testing.T) {
	assert.True(t, (&LsblkDevice{Name: "/dev/zram0"}).IsZram())
}

func TestLsblkDevice_IsZram_Not(t *testing.T) {
	assert.False(t, (&LsblkDevice{Name: "/dev/sda"}).IsZram())
}
