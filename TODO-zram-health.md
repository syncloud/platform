# TODO: zram swap + Health-page bottleneck visibility

## zram swap on high-RAM-but-overloaded boxes (boris prod)

### Context
- boris.syncloud.it prod: 8 GB RAM (`MemTotal 8052952 kB`), ~15 snaps running (photoprism, mastodon, matrix, nextcloud, pihole, games, mail…).
- Only swap = legacy 953 MB disk partition `/dev/sda1`, 100% full (cold parked pages).
- Box was imaged from the **legacy amd64 image** (`debian-buster-amd64-8gb.img`); swap partition + fstab entry come from the image, NOT added manually.
- Modern amd64 pipeline (`image/tools/extract-amd64.sh` → debian-12-generic + EFI, `rootfs-amd64.sh`) ships **no swap** at all. arm/arm64 fstabs also have no swap.

### Why platform zram is skipped here
`backend/stability/zram.go` `EnsureConfigured()`:
```go
memThresholdKB = 6 * 1024 * 1024  // 6 GiB
if snap.TotalKB > memThresholdKB { skip }
```
Box has 8 GB > 6 GiB gate ⇒ zram never runs. Confirmed on prod: zram module not loaded, no `/dev/zram0`.

### Decision: gate removed — zram always enabled (DONE)
`EnsureConfigured()` no longer checks memory; `memThresholdKB` is gone. zram is
configured on every box regardless of total/available RAM. Compressed swap is
cheap headroom even on large boxes, and `alreadyOn()` keeps it idempotent.
(Considered but rejected: Option A bump to 8 GiB; Option B gate on available RAM.)

### Existing swap NOT auto-disabled — no extra code needed
- `disableFileSwaps()` only turns off `Type == "file"` swaps (`fields[1] != "file"` → skip).
- boris `/dev/sda1` is `Type == partition` ⇒ left alone.
- Result after enabling zram: zram prio 10 (used first) + sda1 prio -2 (overflow). Ideal. fstab re-adds sda1 on boot, stability service re-adds zram on boot. Stable, zero further changes.
- Only if we wanted to KILL sda1 too would we need code (extend disableFileSwaps to partitions — risky, affects all devices) or edit this box's fstab + `swapoff`.
- zram would size to 2 GiB cap (RAM/2=4G clamped), zstd comp_algorithm, prio 10.

## Prove disk/swap is (or isn't) the bottleneck + surface it on Health page

### Current evidence: NOT disk-bound (correction to earlier swap-thrash hypothesis)
Sampled prod ~12s + swap counters:
- `vmstat` si/so = **0/0** sustained → not paging in/out of swap.
- wa (iowait) = **0–1%**, idle 95–97%, load 0.2–0.6.
- swap 100% full but **static** (cold pages) — benign, normal.
- cumulative pswpin 92 days = 354k pages ≈ 1.4 GB (~15 MB/day) — negligible.
- The 100% swap BAR is misleading; RATE (si/so) + iowait are the real proof, both ~0 at rest.
- k9mail slowness more likely IMAP-over-TLS RTT or one-off Dovecot indexing. TODO: sample **during** a k9mail sync (si/so, wa, `/proc/<dovecot-pid>/io`).

### Constraint
- PSI unavailable: kernel `4.19.0-8-amd64` (legacy buster), needs ≥4.20 + CONFIG_PSI. `/proc/pressure/*` absent.
- Note: `backend/stability/events.go` already has a `PSIavg10` field → newer kernels/boards would populate it, this box can't.

### Health page work (`backend/health/metrics.go` + `web/platform/src/views/Health.vue`)
Backend already collects: `CPU.IOWait` (unused), `Disk.SectorsRead/Wrt` (shown as ↓/↑ KB/s), `Memory.SwapTotalKB/FreeKB` (swap bar).

Add, in priority order:
1. **Swap-in/out rate (si/so)** — the definitive thrash signal, NOT collected. Add `swap_in_pages`/`swap_out_pages` (from `/proc/vmstat` pswpin/pswpout) to `Memory` struct; frontend computes delta via existing `prevMetrics` pattern (like diskRates) → show `swap ↓X · ↑Y KB/s`.
2. **IOWait %** — already in struct, just not shown. Compute `(iowait_delta/total_delta)*100`, render as its own bar or a segment of the CPU bar.
3. (optional) **PSI mem/io row** — degrade gracefully: hide when `/proc/pressure` absent (this box), show on newer boards.

Goal: turn the misleading static "swap 100%" into signals that distinguish cold-parked-pages (harmless, current state) from active thrashing (real bottleneck).
