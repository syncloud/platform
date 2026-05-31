package support

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/log"
	"strings"
	"testing"
)

func TestLogAggregator_GetLogs(t *testing.T) {
	aggregator := NewAggregator(log.Default())
	logs := aggregator.GetLogs()
	assert.NotEmpty(t, logs)
}

func TestSplitNoise_SeparatesDhcpRegardlessOfTag(t *testing.T) {
	input := "May 22 20:51:49 syncloud dhclient[2301]: XMT: Solicit on eth0\n" +
		"May 22 20:51:49 syncloud sh[2611]: RCV: Advertise message on eth0\n" +
		"May 22 20:51:49 syncloud sh[2611]: RCV:  | X-- t2 - rebind +0\n" +
		"May 22 20:51:49 syncloud sh[2611]: PRC: Lease failed to satisfy.\n" +
		"May 22 20:52:00 syncloud snapd[123]: snap install ok\n" +
		"May 22 20:53:00 syncloud kernel: usb device connected\n"
	signal, noisy := splitNoise(input)
	assert.NotContains(t, signal, "dhclient[")
	assert.NotContains(t, signal, "Advertise")
	assert.NotContains(t, signal, "rebind")
	assert.NotContains(t, signal, "Lease failed")
	assert.Contains(t, signal, "snapd[123]")
	assert.Contains(t, signal, "kernel: usb device connected")
	assert.Contains(t, noisy, "Lease failed")
	assert.Contains(t, noisy, "Advertise")
	assert.Contains(t, noisy, "rebind")
	assert.NotContains(t, noisy, "snapd[123]")
}

func TestSplitNoise_EmptyInput(t *testing.T) {
	signal, noisy := splitNoise("")
	assert.Equal(t, "", signal)
	assert.Equal(t, "", noisy)
}

func TestTail_KeepsLastN(t *testing.T) {
	var lines []string
	for i := 0; i < 1500; i++ {
		lines = append(lines, fmt.Sprintf("line %d", i))
	}
	result := strings.Split(tail(strings.Join(lines, "\n"), 1000), "\n")
	assert.Equal(t, 1000, len(result))
	assert.Equal(t, "line 500", result[0])
	assert.Equal(t, "line 1499", result[999])
}

func TestTail_NoTruncationWhenUnderLimit(t *testing.T) {
	assert.Equal(t, "a\nb\nc", tail("a\nb\nc", 1000))
}
