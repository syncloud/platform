package kernel

import (
	"github.com/stretchr/testify/assert"
	"github.com/syncloud/platform/storage/model"
	"testing"
)

const dmesg = `2026-10-29T10:00:01,000000+00:00 ata1: SATA link up 6.0 Gbps (SStatus 133 SControl 300)
2026-10-29T10:00:01,100000+00:00 ata2: SATA link up 6.0 Gbps (SStatus 133 SControl 300)
2026-10-30T14:12:09,511234+00:00 ata2.00: exception Emask 0x10 SAct 0x0 SErr 0x4050000 action 0xe frozen
2026-10-30T14:12:09,511301+00:00 ata2: SError: { PHYRdyChg CommWake DevExch }
2026-10-30T14:12:09,511344+00:00 ata2.00: failed command: READ FPDMA QUEUED
2026-10-30T14:12:09,511402+00:00 ata2: hard resetting link
2026-10-30T14:12:15,812000+00:00 ata2: SATA link up 6.0 Gbps (SStatus 133 SControl 300)
2026-10-31T02:01:44,000100+00:00 ata2: hard resetting link
2026-10-31T02:01:45,000100+00:00 blk_update_request: I/O error, dev sdb, sector 123456 op 0x0:(READ) flags 0x0 phys_seg 1 prio class 0
2026-10-31T02:01:45,000200+00:00 Buffer I/O error on dev sdb1, logical block 15432, async page read
2026-10-31T02:01:45,000300+00:00 BTRFS error (device sda): bdev /dev/sdb errs: wr 0, rd 1, flush 0, corrupt 0, gen 0
2026-10-31T02:01:46,000000+00:00 I/O error, dev loop10, sector 8 op 0x0:(READ) flags 0x80700 phys_seg 1 prio class 2
2026-10-31T02:01:47,000000+00:00 critical medium error, dev nvme0n1p2, sector 900 op 0x0:(READ)
`

func TestEvents_Link(t *testing.T) {
	events := NewEvents(ParseLines(dmesg))
	assert.Equal(t, model.KernelEvents{Count: 4, Last: "2026-10-31T02:01:44"}, events.Link("ata2"))
	assert.Equal(t, model.KernelEvents{}, events.Link("ata1"))
	assert.Equal(t, model.KernelEvents{}, events.Link(""))
}

func TestEvents_Io(t *testing.T) {
	events := NewEvents(ParseLines(dmesg))
	assert.Equal(t, model.KernelEvents{Count: 3, Last: "2026-10-31T02:01:45"}, events.Io("/dev/sdb"))
	assert.Equal(t, model.KernelEvents{}, events.Io("/dev/sda"))
	assert.Equal(t, model.KernelEvents{}, events.Io("/dev/loop1"))
	assert.Equal(t, model.KernelEvents{Count: 1, Last: "2026-10-31T02:01:46"}, events.Io("/dev/loop10"))
	assert.Equal(t, model.KernelEvents{Count: 1, Last: "2026-10-31T02:01:47"}, events.Io("/dev/nvme0n1"))
}

func TestParseLines_SkipsLinesWithoutTimestamp(t *testing.T) {
	lines := ParseLines("garbage\n2026-10-29T10:00:01,000000+00:00 ata1: hard resetting link\n")
	assert.Equal(t, []Line{{Time: "2026-10-29T10:00:01", Text: "ata1: hard resetting link"}}, lines)
}
