import { installStaleAssetsReload, clearStaleAssetsReload, isStaleAssetError, RELOADED_KEY } from '../../src/util/staleAssets'

function fakeWindow () {
  const store = {}
  const listeners = {}
  return {
    reloads: 0,
    sessionStorage: {
      getItem: (k) => (k in store ? store[k] : null),
      setItem: (k, v) => { store[k] = v },
      removeItem: (k) => { delete store[k] }
    },
    location: { reload () { this.reloads += 1 } },
    addEventListener: (name, handler) => { listeners[name] = handler },
    fire: (name, event) => listeners[name](event)
  }
}

function fakeRouter () {
  let handler = null
  return {
    onError: (h) => { handler = h },
    fail: (error) => handler(error)
  }
}

test('detects stale chunk errors', () => {
  expect(isStaleAssetError(new Error('Failed to fetch dynamically imported module: /assets/App.abc.js'))).toBe(true)
  expect(isStaleAssetError(new Error('Unable to preload CSS for /assets/App.abc.css'))).toBe(true)
  expect(isStaleAssetError(new Error('Request failed with status code 500'))).toBe(false)
})

test('reloads once on router chunk failure', () => {
  const window = fakeWindow()
  window.location.reloads = 0
  const router = fakeRouter()
  installStaleAssetsReload(window, router)

  router.fail(new Error('Failed to fetch dynamically imported module: /assets/App.abc.js'))
  expect(window.location.reloads).toBe(1)

  router.fail(new Error('Failed to fetch dynamically imported module: /assets/App.abc.js'))
  expect(window.location.reloads).toBe(1)
})

test('ignores unrelated router errors', () => {
  const window = fakeWindow()
  window.location.reloads = 0
  const router = fakeRouter()
  installStaleAssetsReload(window, router)

  router.fail(new Error('Request failed with status code 500'))
  expect(window.location.reloads).toBe(0)
})

test('reloads on vite preload error', () => {
  const window = fakeWindow()
  window.location.reloads = 0
  const router = fakeRouter()
  installStaleAssetsReload(window, router)

  let prevented = false
  window.fire('vite:preloadError', { preventDefault: () => { prevented = true } })

  expect(prevented).toBe(true)
  expect(window.location.reloads).toBe(1)
})

test('clearing the flag allows a later reload', () => {
  const window = fakeWindow()
  window.location.reloads = 0
  const router = fakeRouter()
  installStaleAssetsReload(window, router)

  router.fail(new Error('Failed to fetch dynamically imported module: /assets/App.abc.js'))
  expect(window.sessionStorage.getItem(RELOADED_KEY)).toBe('1')

  clearStaleAssetsReload(window)
  expect(window.sessionStorage.getItem(RELOADED_KEY)).toBe(null)

  router.fail(new Error('Failed to fetch dynamically imported module: /assets/App.abc.js'))
  expect(window.location.reloads).toBe(2)
})
