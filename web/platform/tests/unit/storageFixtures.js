export const arrayDisks = [
  {
    name: 'CT2000BX500SSD1',
    device: '/dev/sda',
    active: true,
    size: '1.8T',
    raid: 'raid1',
    raid_metadata: 'raid1',
    serial: '2112E58EB6C6',
    hctl: '0:0:0:0',
    devid: '1',
    errors: { write: 0, read: 0, flush: 0, corruption: 0, generation: 0 },
    has_errors: false,
    partitions: []
  },
  {
    name: 'External_USB_3.0',
    device: '/dev/sdb',
    active: true,
    size: '1.8T',
    raid: 'raid1',
    raid_metadata: 'raid1',
    serial: '20170617007242F',
    hctl: '1:0:0:0',
    devid: '2',
    errors: { write: 0, read: 12, flush: 0, corruption: 3, generation: 0 },
    has_errors: true,
    partitions: []
  }
]

export const arrayHealth = (overrides) => Object.assign({
  disks: [
    {
      device: '/dev/sda',
      port: 'ata1',
      smart: {
        available: true,
        passed: true,
        reallocated: 0,
        pending: 0,
        offline_uncorrectable: 0,
        crc_errors: 0,
        power_on_hours: 212,
        temperature: 31,
        self_test: { running: false, remaining_percent: 0, last_type: 'Short offline', last_status: 'Completed without error', last_passed: true },
        log: [
          { type: 'Short offline', status: 'Completed without error', passed: true, hours: 212 },
          { type: 'Short offline', status: 'Completed: read failure', passed: false, hours: 40 }
        ]
      },
      link: { count: 0 },
      io: { count: 0 }
    },
    {
      device: '/dev/sdb',
      port: 'ata2',
      smart: { available: false, message: 'Unknown USB bridge', self_test: { running: false, remaining_percent: 0 } },
      link: { count: 4, last: '2026-10-31T02:01:44' },
      io: { count: 3, last: '2026-10-31T02:01:45' }
    }
  ],
  scrub: { state: 'finished', started: 'Sat Nov  7 21:40:51 2026', duration: '1:10:02', errors: 'read=12 csum=3', corrected: 0, uncorrectable: 15 },
  schedule: { enabled: false },
  busy: false
}, overrides)
