package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLsblkEntry_IsMMCBootPartition(t *testing.T) {
	assert.True(t, (&LsblkEntry{Name: "/dev/mmcblk1boot0"}).IsMMCBootPartition())
}

func TestLsblkEntry_IsMMCBootPartition_Not(t *testing.T) {
	assert.False(t, (&LsblkEntry{Name: "/dev/mmcblk1p0"}).IsMMCBootPartition())
}

func TestLsblkEntry_IsZram(t *testing.T) {
	assert.True(t, (&LsblkEntry{Name: "/dev/zram0"}).IsZram())
}

func TestLsblkEntry_IsZram_Not(t *testing.T) {
	assert.False(t, (&LsblkEntry{Name: "/dev/sda"}).IsZram())
}
