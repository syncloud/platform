function value (number) {
  return number === undefined || number === null ? '-' : number
}

export function healthOf (health, disk) {
  return health.disks.find(h => h.device === disk.device)
}

export function portOf (disk, health) {
  const port = health && health.port
  if (port && disk.hctl) {
    return port + ' (' + disk.hctl + ')'
  }
  return port || disk.hctl
}

export function diskStatus (disk, health) {
  if (!health) {
    return ''
  }
  const smart = health.smart
  const testFailed = smart.self_test && smart.self_test.last_passed === false
  if (disk.has_errors || health.io.count > 0 || testFailed || (smart.available && (smart.passed === false || smart.pending > 0 || smart.offline_uncorrectable > 0))) {
    return 'error'
  }
  if (health.link.count > 0 || (smart.available && (smart.crc_errors > 0 || smart.reallocated > 0))) {
    return 'warning'
  }
  return 'ok'
}

export function smartBad (smart) {
  return smart.available && (smart.passed === false || [smart.reallocated, smart.pending, smart.offline_uncorrectable, smart.crc_errors].some(v => v > 0))
}

export function smartText (t, smart) {
  if (!smart.available) {
    return t('storage.smartNotAvailable')
  }
  return t(smart.passed ? 'storage.smartPassed' : 'storage.smartFailed', {
    hours: value(smart.power_on_hours),
    temperature: value(smart.temperature),
    reallocated: value(smart.reallocated),
    pending: value(smart.pending),
    uncorrectable: value(smart.offline_uncorrectable),
    crc: value(smart.crc_errors)
  })
}

export function selfTestText (t, test) {
  if (test.running) {
    return t('storage.selfTestRunning', { percent: test.remaining_percent })
  }
  if (!test.last_status) {
    return t('storage.selfTestNone')
  }
  return [test.last_type, test.last_status].filter(v => v).join(': ')
}

export function selfTestsText (t, health) {
  const disks = health.disks.filter(d => d.smart.available)
  const names = list => list.map(d => d.device).join(', ')
  if (disks.length === 0) {
    return t('storage.selfTestsNotAvailable')
  }
  const running = disks.filter(d => d.smart.self_test.running)
  if (running.length > 0) {
    return t('storage.selfTestsRunning', { disks: names(running) })
  }
  const failed = disks.filter(d => d.smart.self_test.last_passed === false)
  if (failed.length > 0) {
    return t('storage.selfTestsFailed', { disks: names(failed) })
  }
  if (disks.every(d => !d.smart.self_test.last_status)) {
    return t('storage.selfTestNone')
  }
  return t('storage.selfTestsPassed')
}

export function selfTestsBad (health) {
  return health.disks.some(d => d.smart.available && !d.smart.self_test.running && d.smart.self_test.last_passed === false)
}

export function checksRunning (health) {
  return (!!health.scrub && health.scrub.state === 'running') || health.disks.some(d => d.smart.available && d.smart.self_test.running)
}

export function eventsText (t, events) {
  if (events.count === 0) {
    return t('storage.eventsNone')
  }
  return t('storage.eventsCount', { count: events.count, last: events.last })
}

export function scrubText (t, scrub) {
  let text
  if (scrub.state === 'never') {
    text = t('storage.scrubNever')
  } else if (scrub.state === 'running') {
    text = t('storage.scrubRunning', { progress: scrub.progress || '0%' })
  } else if (scrub.state === 'finished') {
    text = t('storage.scrubFinished', { started: scrub.started, duration: scrub.duration })
  } else {
    text = t('storage.scrubStopped', { started: scrub.started })
  }
  if (scrub.errors) {
    return text + '. ' + t('storage.scrubErrors', { errors: scrub.errors, corrected: scrub.corrected, uncorrectable: scrub.uncorrectable })
  }
  return text
}

export function logEntryText (t, entry) {
  return t('storage.selfTestLogEntry', { type: entry.type, status: entry.status, hours: entry.hours })
}
