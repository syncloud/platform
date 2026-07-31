import { Page, expect } from '@playwright/test'
import { waitForLoading } from './ui'

const deviceUser = process.env.PLAYWRIGHT_DEVICE_USER ?? 'user'
const devicePassword = process.env.PLAYWRIGHT_DEVICE_PASSWORD ?? 'Password1'

export { deviceUser, devicePassword, waitForLoading }

const loginAttempts = 3

export async function login(page: Page, opts: { user?: string; password?: string } = {}) {
  const applications = page.getByRole('heading', { name: 'Applications' })
  for (let attempt = 1; attempt <= loginAttempts; attempt++) {
    await page.goto('/')
    if (await applications.isVisible()) {
      break
    }
    const firstFactor = page
      .waitForResponse(response => response.url().includes('/api/firstfactor'), { timeout: 30000 })
      .catch(() => null)
    await page.locator('#username-textfield').fill(opts.user ?? deviceUser)
    await page.locator('#password-textfield').fill(opts.password ?? devicePassword)
    await page.locator('#sign-in-button').click()
    const response = await firstFactor
    if (response !== null && response.status() >= 500 && attempt < loginAttempts) {
      await page.waitForTimeout(2000)
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
