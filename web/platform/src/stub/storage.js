import { Response } from 'miragejs'

const busyMessage = 'an app install or update is in progress, try again when it finishes'

const zeroErrors = () => ({ write: 0, read: 0, flush: 0, corruption: 0, generation: 0 })

const smartOk = (hours, temperature, crc) => ({
  available: true,
  passed: true,
  reallocated: 0,
  pending: 0,
  offline_uncorrectable: 0,
  crc_errors: crc,
  power_on_hours: hours,
  temperature: temperature,
  self_test: { running: false, remaining_percent: 0, last_type: 'Short offline', last_status: 'Completed without error', last_passed: true }
})

const arrayDisks = () => [
  {
    name: 'CT2000BX500SSD1',
    device: '/dev/sda',
    size: '1.8T',
    serial: '2112E58EB6C6',
    hctl: '0:0:0:0',
    port: 'ata1',
    errors: zeroErrors(),
    smart: smartOk(8784, 26, 0),
    link: { count: 0 },
    io: { count: 0 }
  },
  {
    name: 'WDC WD20EFRX-68EUZN0',
    device: '/dev/sdb',
    size: '1.8T',
    serial: 'WD-WCC4M1KX2H7A',
    hctl: '1:0:0:0',
    port: 'ata2',
    errors: { write: 0, read: 12, flush: 0, corruption: 3, generation: 0 },
    smart: smartOk(21430, 34, 41),
    link: { count: 4, last: '2026-10-31T02:01:44' },
    io: { count: 3, last: '2026-10-31T02:01:45' }
  },
  {
    name: 'Seagate Expansion Desk with a rather long product name',
    device: '/dev/sdc',
    size: '3.6T',
    serial: 'NA8F2XKQ',
    hctl: '2:0:0:0',
    port: 'usb2',
    errors: zeroErrors(),
    smart: { available: false, message: 'Unknown USB bridge [0x0bc2:0x3312 (0x100)]', self_test: { running: false, remaining_percent: 0 } },
    link: { count: 0 },
    io: { count: 0 }
  },
  {
    name: 'ST2000DM008-2FR102',
    device: '/dev/sdd',
    size: '1.8T',
    serial: 'ZFL3K9QW',
    hctl: '3:0:0:0',
    port: 'ata4',
    errors: zeroErrors(),
    smart: {
      available: true,
      passed: false,
      reallocated: 120,
      pending: 8,
      offline_uncorrectable: 8,
      crc_errors: 0,
      power_on_hours: 40211,
      temperature: 41,
      self_test: { running: false, remaining_percent: 0 }
    },
    link: { count: 0 },
    io: { count: 0 }
  }
]

export function storageStub (scenario) {
  const inArray = scenario !== 'noarray'
  const disks = scenario === 'empty' ? [] : arrayDisks()
  const state = {
    busy: scenario === 'busy',
    scrub: inArray && disks.length > 0 ? { state: 'never', corrected: 0, uncorrectable: 0 } : null,
    scrubSteps: [],
    schedule: { scrub: false, smart: false }
  }

  const diskList = () => disks.map((disk, index) => {
    const result = {
      name: disk.name,
      device: disk.device,
      size: disk.size,
      serial: disk.serial,
      hctl: disk.hctl,
      active: inArray,
      partitions: [],
      has_errors: false
    }
    if (inArray) {
      result.uuid = 'b7020451-ac9c-4195-a7e8-f072458b1bae'
      result.raid = 'raid1'
      result.raid_metadata = 'raid1c3'
      result.devid = String(index + 1)
      result.errors = { ...disk.errors }
      result.has_errors = Object.values(disk.errors).some(v => v > 0)
    } else {
      result.partitions = [{ size: disk.size, device: disk.device + '1', mount_point: '', active: index === 0, fs_type: 'ext4', extendable: false }]
    }
    return result
  })

  const advance = () => {
    if (state.scrub && state.scrub.state === 'running') {
      const next = state.scrubSteps.shift()
      Object.assign(state.scrub, next)
    }
    disks.forEach(disk => {
      const test = disk.smart.self_test
      if (!test.running) {
        return
      }
      test.remaining_percent -= 30
      if (test.remaining_percent <= 0) {
        const failing = disk.smart.passed === false
        disk.smart.self_test = {
          running: false,
          remaining_percent: 0,
          last_type: 'Short offline',
          last_status: failing ? 'Completed: read failure' : 'Completed without error',
          last_passed: !failing,
          last_hours: disk.smart.power_on_hours
        }
      }
    })
  }

  const refused = () => state.busy ? new Response(500, {}, { success: false, message: busyMessage }) : null

  return function routes () {
    this.get('/rest/storage/disks', function () {
      return new Response(200, {}, { success: true, data: diskList() })
    })
    this.get('/rest/storage/health', function () {
      advance()
      return new Response(200, {}, {
        success: true,
        data: {
          disks: disks.map(disk => ({ device: disk.device, port: disk.port, smart: disk.smart, link: disk.link, io: disk.io })),
          scrub: state.scrub,
          schedule: state.schedule,
          busy: state.busy
        }
      })
    })
    this.post('/rest/storage/check/scrub', function () {
      if (!state.scrub) {
        return new Response(500, {}, { success: false, message: 'no active btrfs disk array' })
      }
      if (state.scrub.state === 'running') {
        return new Response(500, {}, { success: false, message: 'a scrub is already running' })
      }
      const busy = refused()
      if (busy) {
        return busy
      }
      const damaged = disks.some(disk => Object.values(disk.errors).some(v => v > 0))
      const started = new Date().toString().slice(0, 24)
      state.scrub = { state: 'running', started: started, duration: '0:00:00', progress: '0.00%', corrected: 0, uncorrectable: 0 }
      state.scrubSteps = [
        { progress: '25.40%', duration: '0:17:52' },
        { progress: '61.02%', duration: '0:42:10' },
        damaged
          ? { state: 'finished', progress: '', duration: '1:10:02', errors: 'read=12 csum=3', corrected: 3, uncorrectable: 12 }
          : { state: 'finished', progress: '', duration: '1:08:47', errors: '', corrected: 0, uncorrectable: 0 }
      ]
      return new Response(200, {}, { success: true, data: 'OK' })
    })
    this.post('/rest/storage/check/smart', function (_schema, request) {
      const device = JSON.parse(request.requestBody).device
      const disk = disks.find(d => d.device === device)
      if (!disk) {
        return new Response(500, {}, { success: false, message: 'unknown disk: ' + device })
      }
      const busy = refused()
      if (busy) {
        return busy
      }
      if (!disk.smart.available) {
        return new Response(500, {}, { success: false, message: disk.smart.message })
      }
      disk.smart.self_test = { running: true, remaining_percent: 90 }
      return new Response(200, {}, { success: true, data: 'OK' })
    })
    this.get('/rest/storage/check/schedule', function () {
      return new Response(200, {}, { success: true, data: state.schedule })
    })
    this.post('/rest/storage/check/schedule', function (_schema, request) {
      state.schedule = JSON.parse(request.requestBody)
      return new Response(200, {}, { success: true, data: state.schedule })
    })
    this.post('/rest/storage/errors/reset', function () {
      if (!state.scrub) {
        return new Response(500, {}, { success: false, message: 'no active btrfs disk array' })
      }
      disks.forEach(disk => { disk.errors = zeroErrors() })
      return new Response(200, {}, { success: true, data: 'OK' })
    })
  }
}
