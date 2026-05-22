package support

import (
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

func TestFilterAndTail_ExcludesMatchingLines(t *testing.T) {
	input := "May 22 20:51:49 syncloud dhclient[2301]: XMT: Solicit on eth0\n" +
		"May 22 20:51:49 syncloud dhclient[2301]: RCV: Advertise message on eth0\n" +
		"May 22 20:52:00 syncloud snapd[123]: snap install ok\n" +
		"May 22 20:53:00 syncloud kernel: usb device connected\n"
	out := filterAndTail(input, "dhclient[", 1000)
	assert.NotContains(t, out, "dhclient[")
	assert.Contains(t, out, "snapd[123]")
	assert.Contains(t, out, "kernel: usb device connected")
}

func TestFilterAndTail_KeepsLastNAfterFilter(t *testing.T) {
	var lines []string
	for i := 0; i < 1500; i++ {
		if i%3 == 0 {
			lines = append(lines, "dhclient[1]: noise")
		} else {
			lines = append(lines, "keeper line")
		}
	}
	out := filterAndTail(strings.Join(lines, "\n"), "dhclient[", 1000)
	resultLines := strings.Split(out, "\n")
	assert.Equal(t, 1000, len(resultLines))
	for _, l := range resultLines {
		assert.Equal(t, "keeper line", l)
	}
}

func TestFilterAndTail_NoTruncationWhenUnderLimit(t *testing.T) {
	input := "a\nb\nc"
	out := filterAndTail(input, "x", 1000)
	assert.Equal(t, "a\nb\nc", out)
}

func TestFilterAndTail_EmptyInput(t *testing.T) {
	out := filterAndTail("", "dhclient[", 1000)
	assert.Equal(t, "", out)
}
