function createMask (text) {
  const mask = document.createElement('div')
  mask.className = 'sc-loading-mask'
  mask.setAttribute('data-testid', 'loading-mask')
  const spinner = document.createElement('div')
  spinner.className = 'sc-loading-spinner'
  mask.appendChild(spinner)
  if (text) {
    const label = document.createElement('div')
    label.className = 'sc-loading-text'
    label.textContent = text
    mask.appendChild(label)
  }
  return mask
}

const Loading = {
  service (options = {}) {
    const mask = createMask(options.text)
    document.body.appendChild(mask)
    let closed = false
    return {
      close () {
        if (closed) return
        closed = true
        if (mask.parentNode) mask.parentNode.removeChild(mask)
      }
    }
  }
}

export default Loading
