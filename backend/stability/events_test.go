package stability

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendAndRecentReverseOrder(t *testing.T) {
	dir := t.TempDir()
	log := NewEventLog(filepath.Join(dir, "events.jsonl"))
	require.NoError(t, log.Append(Event{Kind: EventKindZramEnabled, SizeBytes: 1 << 30}))
	require.NoError(t, log.Append(Event{Kind: EventKindPressure, AvailRatio: 0.05}))
	require.NoError(t, log.Append(Event{Kind: EventKindVictimSigterm, PID: 1234, Comm: "python3", RSSkb: 2000000}))

	evs, err := log.Recent(10)
	require.NoError(t, err)
	require.Len(t, evs, 3)
	assert.Equal(t, EventKindVictimSigterm, evs[0].Kind)
	assert.Equal(t, "python3", evs[0].Comm)
	assert.Equal(t, EventKindPressure, evs[1].Kind)
	assert.Equal(t, EventKindZramEnabled, evs[2].Kind)
}

func TestRecentMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	evs, err := NewEventLog(filepath.Join(dir, "nope.jsonl")).Recent(10)
	require.NoError(t, err)
	assert.Empty(t, evs)
}

func TestRecentCapsLimit(t *testing.T) {
	dir := t.TempDir()
	log := NewEventLog(filepath.Join(dir, "events.jsonl"))
	for i := 0; i < 20; i++ {
		require.NoError(t, log.Append(Event{Kind: EventKindPressure, PID: i}))
	}
	evs, err := log.Recent(5)
	require.NoError(t, err)
	require.Len(t, evs, 5)
	assert.Equal(t, 19, evs[0].PID)
	assert.Equal(t, 15, evs[4].PID)
}
