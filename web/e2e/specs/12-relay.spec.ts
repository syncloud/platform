import { test, expect, Page } from '@playwright/test'
import { execFileSync } from 'node:child_process'
import { login } from '../helpers/login'
import { settings } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'
import { ssh } from '../helpers/ssh'

test.describe.configure({ mode: 'serial' })

const domain = process.env.PLAYWRIGHT_DOMAIN ?? 'bookworm-amd64.redirect'
const relayHost = 'relay.redirect'

function httpCode(args: string[]): string {
  try {
    return execFileSync('curl', ['-sk', '-o', '/dev/null', '-w', '%{http_code}', '--max-time', '15', ...args],
      { encoding: 'utf8' }).trim()
  } catch (e: any) {
    return '000'
  }
}

function throughRelay(): string {
  return httpCode(['--connect-to', `${domain}:443:${relayHost}:4443`, `https://${domain}/`])
}

function direct(): string {
  return httpCode([`https://${domain}/`])
}

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  await login(page)
})

test.afterAll(async () => {
  // saving access fires the access-change app fan-out in the background; let it
  // drain so its OIDC/nginx reconfigure doesn't leak into the next spec's login
  for (let i = 0; i < 60; i++) {
    const running = ssh("pgrep -f 'access[-]change' || true", { throw: false })
    if (!running.trim()) break
    await new Promise((r) => setTimeout(r, 1000))
  }
  await page.close()
})

async function save(testInfo: any, name: string) {
  const mask = page.getByTestId('loading-mask')
  await page.getByTestId('access-save').click()
  await expect(mask).toBeVisible()
  await expect(mask).toBeHidden({ timeout: 40000 })
  await shoot(page, testInfo, name)
  await expect(page.getByTestId('error-message')).toBeHidden()
}

test('relay enable saves without error', async ({}, testInfo) => {
  await settings(page, 'access', testInfo)
  await expect(page.getByTestId('ipv4-section')).toBeVisible()
  await page.getByTestId('ipv4-mode-relay').click()
  await save(testInfo, 'relay_enabled')
  await expect(page.getByTestId('relay-status')).toBeVisible()
})

test('traffic tunnels through frps when relay is on', async ({}) => {
  await expect.poll(throughRelay, { timeout: 30000 }).toMatch(/^(200|301|302|401)$/)
})

test('relay disable saves without error', async ({}, testInfo) => {
  await page.getByTestId('ipv4-mode-off').click()
  await save(testInfo, 'relay_disabled')
})

test('traffic is direct and no longer tunnels when relay is off', async ({}) => {
  await expect.poll(direct, { timeout: 30000 }).toMatch(/^(200|301|302|401)$/)
  await expect.poll(throughRelay, { timeout: 30000 }).not.toMatch(/^(200|301|302|401)$/)
})
