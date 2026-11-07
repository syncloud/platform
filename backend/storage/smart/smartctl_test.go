package smart

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

func TestHealth_AtaWithoutSelfTestLog(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{"-j -a /dev/sda": "ata-no-self-test-log.json"}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/sda")

	assert.True(t, health.Available)
	assert.True(t, *health.Passed)
	assert.Nil(t, health.Reallocated)
	assert.Nil(t, health.Pending)
	assert.Equal(t, int64(0), *health.CrcErrors)
	assert.Equal(t, int64(8784), *health.PowerOnHours)
	assert.Equal(t, int64(26), *health.Temperature)
	assert.False(t, health.SelfTest.Running)
	assert.Equal(t, "completed without error", health.SelfTest.LastStatus)
	assert.True(t, *health.SelfTest.LastPassed)
	assert.Equal(t, []string{"-j -a /dev/sda"}, executor.calls)
}

func TestHealth_SelfTestRunning(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{"-j -a /dev/sdb": "ata-self-test-running.json"}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/sdb")

	assert.True(t, health.Available)
	assert.Equal(t, int64(0), *health.Reallocated)
	assert.Equal(t, int64(0), *health.Pending)
	assert.Equal(t, int64(0), *health.OfflineUncorrectable)
	assert.Equal(t, int64(0), *health.CrcErrors)
	assert.True(t, health.SelfTest.Running)
	assert.Equal(t, 90, health.SelfTest.RemainingPercent)
}

func TestHealth_SelfTestDone(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{"-j -a /dev/sdb": "ata-self-test-done.json"}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/sdb")

	assert.False(t, health.SelfTest.Running)
	assert.Equal(t, "Short offline", health.SelfTest.LastType)
	assert.Equal(t, "Completed without error", health.SelfTest.LastStatus)
	assert.True(t, *health.SelfTest.LastPassed)
	assert.Equal(t, int64(213), *health.SelfTest.LastHours)
	passed := true
	assert.Equal(t, []model.SmartLogEntry{{Type: "Short offline", Status: "Completed without error", Passed: &passed, Hours: 213}}, health.Log)
}

func TestHealth_NotAvailable(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{
		"-j -a /dev/loop0":        "unknown-device.json",
		"-j -d sat -a /dev/loop0": "unknown-device.json",
	}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/loop0")

	assert.False(t, health.Available)
	assert.Equal(t, "/dev/loop0: Unable to detect device type", health.Message)
	assert.Nil(t, health.Passed)
	assert.Equal(t, []string{"-j -a /dev/loop0", "-j -d sat -a /dev/loop0"}, executor.calls)
}

func TestHealth_UsbBridgeNeedsSat(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{
		"-j -a /dev/sdc":        "unknown-device.json",
		"-j -d sat -a /dev/sdc": "ata-self-test-done.json",
	}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/sdc")

	assert.True(t, health.Available)
	assert.True(t, *health.Passed)
}

func TestHealth_SmartctlMissing(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{}}
	health := NewSmartctl(executor, log.Default()).Health("/dev/sda")

	assert.False(t, health.Available)
	assert.Equal(t, "exit status 127", health.Message)
}

func TestStartShortTest_UsesDetectedDeviceType(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{
		"-j -a /dev/sdc":              "unknown-device.json",
		"-j -d sat -a /dev/sdc":       "ata-self-test-done.json",
		"-j -d sat -t short /dev/sdc": "ata-self-test-done.json",
	}}
	err := NewSmartctl(executor, log.Default()).StartShortTest("/dev/sdc")

	assert.Nil(t, err)
	assert.Equal(t, "-j -d sat -t short /dev/sdc", executor.calls[2])
}

func TestStartShortTest_NotAvailable(t *testing.T) {
	executor := &ExecutorStub{responses: map[string]string{
		"-j -a /dev/loop0":        "unknown-device.json",
		"-j -d sat -a /dev/loop0": "unknown-device.json",
	}}
	err := NewSmartctl(executor, log.Default()).StartShortTest("/dev/loop0")

	assert.EqualError(t, err, "/dev/loop0: Unable to detect device type")
	assert.Len(t, executor.calls, 2)
}
