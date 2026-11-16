const browsers = [
  [/Edg\//, 'Edge'],
  [/Firefox\//, 'Firefox'],
  [/Chrome\//, 'Chrome'],
  [/Safari\//, 'Safari']
]

const systems = [
  [/Android/, 'Android'],
  [/iPhone|iPad/, 'iOS'],
  [/Windows/, 'Windows'],
  [/Mac OS X/, 'macOS'],
  [/Linux/, 'Linux']
]

function match (list, ua, fallback) {
  const found = list.find(([pattern]) => pattern.test(ua))
  return found ? found[1] : fallback
}

export function deviceName (ua = navigator.userAgent) {
  return match(browsers, ua, 'Browser') + ' on ' + match(systems, ua, 'unknown device')
}
