package btrfs

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

func TestParseScrubStatus_Never(t *testing.T) {
	status := ParseScrubStatus(`UUID:             a13d8208-9d05-4587-98a4-c790cbcf93e0
	no stats available
Total to scrub:   100.41MiB
Rate:             0.00B/s
Error summary:    no errors found
`)
	assert.Equal(t, model.ScrubStatus{State: model.ScrubNever}, status)
}

func TestParseScrubStatus_JustStarted(t *testing.T) {
	status := ParseScrubStatus(`UUID:             b7020451-ac9c-4195-a7e8-f072458b1bae
	no stats available
Time left:        0:00:00
ETA:              Sat Nov  7 21:40:51 2026
Total to scrub:   160.48MiB
Bytes scrubbed:   0.00B  (0.00%)
Rate:             0.00B/s
Error summary:    no errors found
`)
	assert.Equal(t, model.ScrubRunning, status.State)
	assert.Equal(t, "0.00%", status.Progress)
}

func TestParseScrubStatus_Running(t *testing.T) {
	status := ParseScrubStatus(`UUID:             b7020451-ac9c-4195-a7e8-f072458b1bae
Scrub started:    Sat Nov  7 21:40:51 2026
Status:           running
Duration:         0:12:03
Time left:        1:02:10
ETA:              Sat Nov  7 22:55:04 2026
Total to scrub:   1.20TiB
Bytes scrubbed:   201.33GiB  (16.32%)
Rate:             285.03MiB/s
Error summary:    no errors found
`)
	assert.Equal(t, model.ScrubRunning, status.State)
	assert.Equal(t, "Sat Nov  7 21:40:51 2026", status.Started)
	assert.Equal(t, "0:12:03", status.Duration)
	assert.Equal(t, "16.32%", status.Progress)
	assert.Equal(t, "", status.Errors)
}

func TestParseScrubStatus_Finished(t *testing.T) {
	status := ParseScrubStatus(`UUID:             b7020451-ac9c-4195-a7e8-f072458b1bae
Scrub started:    Sat Nov  7 21:40:51 2026
Status:           finished
Duration:         0:00:00
Total to scrub:   160.48MiB
Rate:             0.00B/s
Error summary:    no errors found
`)
	assert.Equal(t, model.ScrubStatus{State: model.ScrubFinished, Started: "Sat Nov  7 21:40:51 2026", Duration: "0:00:00"}, status)
}

func TestParseScrubStatus_Aborted(t *testing.T) {
	status := ParseScrubStatus(`UUID:             b7020451-ac9c-4195-a7e8-f072458b1bae
Scrub started:    Sat Nov  7 21:40:57 2026
Status:           aborted
Duration:         0:00:00
Total to scrub:   160.48MiB
Rate:             0.00B/s
Error summary:    no errors found
`)
	assert.Equal(t, model.ScrubAborted, status.State)
}

func TestParseScrubStatus_Errors(t *testing.T) {
	status := ParseScrubStatus(`UUID:             b7020451-ac9c-4195-a7e8-f072458b1bae
Scrub started:    Sat Nov  7 21:40:51 2026
Status:           finished
Duration:         1:10:02
Total to scrub:   1.20TiB
Rate:             298.12MiB/s
Error summary:    read=12 csum=3
  Corrected:      0
  Uncorrectable:  15
  Unverified:     0
`)
	assert.Equal(t, model.ScrubFinished, status.State)
	assert.Equal(t, "read=12 csum=3", status.Errors)
	assert.Equal(t, int64(0), status.Corrected)
	assert.Equal(t, int64(15), status.Uncorrectable)
}
