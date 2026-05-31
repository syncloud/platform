import fs from 'fs'
import path from 'path'

const localesDir = path.resolve(__dirname, '../../src/locales')
const reference = 'en'

function keyPaths (obj, prefix = '') {
  const paths = []
  for (const [key, value] of Object.entries(obj)) {
    const full = prefix ? `${prefix}.${key}` : key
    if (value !== null && typeof value === 'object' && !Array.isArray(value)) {
      paths.push(...keyPaths(value, full))
    } else {
      paths.push(full)
    }
  }
  return paths
}

function load (locale) {
  return JSON.parse(fs.readFileSync(path.join(localesDir, `${locale}.json`), 'utf8'))
}

const referenceKeys = keyPaths(load(reference))
const locales = fs.readdirSync(localesDir)
  .filter(f => f.endsWith('.json'))
  .map(f => f.replace(/\.json$/, ''))
  .filter(l => l !== reference)

test('all locales are present', () => {
  expect(locales.length).toBeGreaterThan(0)
})

test.each(locales)('locale %s has no key gaps against en', (locale) => {
  const keys = keyPaths(load(locale))
  const missing = referenceKeys.filter(k => !keys.includes(k))
  const extra = keys.filter(k => !referenceKeys.includes(k))
  expect({ locale, missing, extra }).toEqual({ locale, missing: [], extra: [] })
})
