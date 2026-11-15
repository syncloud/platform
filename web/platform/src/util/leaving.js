let leaving = false

window.addEventListener('pageshow', () => {
  leaving = false
})

export function leave (url) {
  leaving = true
  window.location.href = url
}

export function isLeaving () {
  return leaving
}
