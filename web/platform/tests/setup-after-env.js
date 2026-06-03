jest.setTimeout(30000)

// Fail tests on unexpected console output so warnings/errors cannot slip
// through a green build. Known-benign messages are suppressed below.
const IGNORE = [
  // <Transition @after-enter> warns in the jsdom/test-utils environment even
  // though the handler is a function; it works at runtime and the data-ready
  // hook it sets is required by the access-page e2e screenshot timing.
  /Wrong type passed as event handler to onAfterEnter/
]

const originalError = console.error
const originalWarn = console.warn
let captured = []

function capture (kind, original) {
  return (...args) => {
    const message = args.map(a => (a && a.stack) ? a.stack : String(a)).join(' ')
    if (IGNORE.some(re => re.test(message))) return
    captured.push(kind + ': ' + message)
    original(...args)
  }
}

beforeEach(() => { captured = [] })

afterEach(() => {
  if (captured.length > 0) {
    const lines = captured.join('\n')
    captured = []
    throw new Error('Unexpected console output during test:\n' + lines)
  }
})

console.error = capture('console.error', originalError)
console.warn = capture('console.warn', originalWarn)
