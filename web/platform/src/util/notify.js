function ensureContainer () {
  let c = document.querySelector('.sc-toast-container')
  if (!c) {
    c = document.createElement('div')
    c.className = 'sc-toast-container'
    document.body.appendChild(c)
  }
  return c
}

export default function notify (options = {}) {
  const container = ensureContainer()
  const el = document.createElement('div')
  el.className = 'sc-toast sc-toast-' + (options.type || 'info')
  el.setAttribute('data-testid', 'toast')

  if (options.title) {
    const title = document.createElement('div')
    title.className = 'sc-toast-title'
    title.textContent = options.title
    el.appendChild(title)
  }
  if (options.message) {
    const msg = document.createElement('div')
    msg.className = 'sc-toast-message'
    msg.textContent = options.message
    el.appendChild(msg)
  }
  const closeBtn = document.createElement('button')
  closeBtn.className = 'sc-toast-close'
  closeBtn.type = 'button'
  closeBtn.setAttribute('aria-label', 'close')
  closeBtn.textContent = '✕'
  el.appendChild(closeBtn)

  container.appendChild(el)

  let removed = false
  const remove = () => {
    if (removed) return
    removed = true
    if (el.parentNode) el.parentNode.removeChild(el)
    if (typeof options.onClose === 'function') options.onClose()
  }
  closeBtn.addEventListener('click', remove)
  const duration = options.duration === undefined ? 4500 : options.duration
  if (duration > 0) {
    setTimeout(remove, duration)
  }
  return { close: remove }
}
