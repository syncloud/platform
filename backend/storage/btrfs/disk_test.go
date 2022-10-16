package btrfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DetectChange_Replaced(t *testing.T) {
	change, err := DetectChange([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1", "/dev/loop3"})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"replace /dev/loop2 with /dev/loop3"}, change)
}

func Test_DetectChange_Add_One(t *testing.T) {
	change, err := DetectChange([]string{"/dev/loop1"}, []string{"/dev/loop1", "/dev/loop2"})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"add /dev/loop2"}, change)
}

func Test_DetectChange_Create_One(t *testing.T) {
	change, err := DetectChange([]string{}, []string{"/dev/loop1"})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"create /dev/loop1"}, change)
}

func Test_DetectChange_Create_Two(t *testing.T) {
	change, err := DetectChange([]string{}, []string{"/dev/loop1", "/dev/loop2"})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"create /dev/loop1 /dev/loop2"}, change)
}

func Test_DetectChange_Remove_One(t *testing.T) {
	change, err := DetectChange([]string{"/dev/loop1", "/dev/loop2"}, []string{"/dev/loop1"})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"remove /dev/loop2"}, change)
}

func Test_DetectChange_Disable_One(t *testing.T) {
	change, err := DetectChange([]string{"/dev/loop1"}, []string{})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"disable /dev/loop1"}, change)
}

func Test_DetectChange_Disable_Two(t *testing.T) {
	change, err := DetectChange([]string{"/dev/loop1", "/dev/loop2"}, []string{})

	assert.Nil(t, err)
	assert.Len(t, change, 1)
	assert.Equal(t, []string{"disable /dev/loop1 /dev/loop2"}, change)
}
