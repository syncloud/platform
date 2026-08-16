const RELOADED_KEY = 'syncloud-stale-assets-reloaded'

function isStaleAssetError (error) {
  const message = (error && error.message) || String(error || '')
  return message.includes('Failed to fetch dynamically imported module') ||
    message.includes('error loading dynamically imported module') ||
    message.includes('Importing a module script failed') ||
    message.includes('Unable to preload CSS')
}

function reloadOnce (storage, location) {
  if (storage.getItem(RELOADED_KEY)) {
    return false
  }
  storage.setItem(RELOADED_KEY, '1')
  location.reload()
  return true
}

export function installStaleAssetsReload (window, router) {
  const storage = window.sessionStorage

  window.addEventListener('vite:preloadError', (event) => {
    event.preventDefault()
    reloadOnce(storage, window.location)
  })

  router.onError((error) => {
    if (isStaleAssetError(error)) {
      reloadOnce(storage, window.location)
    }
  })
}

export function clearStaleAssetsReload (window) {
  window.sessionStorage.removeItem(RELOADED_KEY)
}

export { isStaleAssetError, RELOADED_KEY }
