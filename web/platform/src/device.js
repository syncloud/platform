import axios from 'axios'

export function deviceHost (deviceUrl) {
  try {
    return new URL(deviceUrl).host
  } catch {
    return ''
  }
}

export function onDeviceUrl (deviceUrl) {
  const host = deviceHost(deviceUrl)
  return host === '' || host === window.location.host
}

export function deviceReachable (deviceUrl, timeout) {
  return axios.get(deviceUrl + '/ping', { timeout: timeout || 5000 })
    .then(() => true)
    .catch(() => false)
}

export function waitForDevice (deviceUrl, timeoutMs, intervalMs) {
  const interval = intervalMs || 2000
  const deadline = Date.now() + (timeoutMs || 120000)
  const attempt = () => deviceReachable(deviceUrl, interval).then(reachable => {
    if (reachable) {
      return true
    }
    if (Date.now() >= deadline) {
      return false
    }
    return new Promise(resolve => setTimeout(resolve, interval)).then(attempt)
  })
  return attempt()
}

export function moveToDevice (deviceUrl) {
  if (onDeviceUrl(deviceUrl)) {
    return Promise.resolve(false)
  }
  return deviceReachable(deviceUrl).then(reachable => {
    if (!reachable) {
      return false
    }
    window.location.replace(deviceUrl + window.location.pathname + window.location.search)
    return true
  })
}
