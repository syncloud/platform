import { Page, expect } from '@playwright/test'
import { waitForLoading } from './ui'

const deviceUser = process.env.PLAYWRIGHT_DEVICE_USER ?? 'user'
const devicePassword = process.env.PLAYWRIGHT_DEVICE_PASSWORD ?? 'Password1'

export { deviceUser, devicePassword, waitForLoading }

const loginAttempts = 2
const formTimeout = 5000
const firstFactorTimeout = 10000

export async function login(page: Page, opts: { user?: string; password?: string } = {}) {
  const applications = page.getByRole('heading', { name: 'Applications' })
  const username = page.locator('#username-textfield')
  for (let attempt = 1; attempt <= loginAttempts; attempt++) {
    const last = attempt === loginAttempts
    await page.goto('/')
    if (await applications.isVisible()) {
      break
    }
    try {
      await expect(username).toBeVisible({ timeout: formTimeout })
    } catch (e) {
      if (last) {
        throw e
      }
      continue
    }
    const firstFactor = page
      .waitForResponse(response => response.url().includes('/api/firstfactor'), { timeout: firstFactorTimeout })
      .catch(() => null)
    await username.fill(opts.user ?? deviceUser)
    await page.locator('#password-textfield').fill(opts.password ?? devicePassword)
    await page.locator('#sign-in-button').click()
    const response = await firstFactor
    if (!last && response !== null && response.status() >= 500) {
      continue
    }
    await expect(applications).toBeVisible()
    break
  }
  await waitForLoading(page)
}

export async function logout(page: Page) {
  const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''
  const url = fullDomain ? `https://${fullDomain}/rest/logout` : '/rest/logout'
  await page.goto(url)
  await expect(page.locator('#username-textfield')).toBeVisible()
}
