import { Page, expect } from '@playwright/test'
import { waitForLoading } from './ui'

const deviceUser = process.env.PLAYWRIGHT_DEVICE_USER ?? 'user'
const devicePassword = process.env.PLAYWRIGHT_DEVICE_PASSWORD ?? 'Password1'

export { deviceUser, devicePassword, waitForLoading }

export async function login(page: Page, opts: { user?: string; password?: string } = {}) {
  await page.goto('/')
  await page.locator('#username-textfield').fill(opts.user ?? deviceUser)
  await page.locator('#password-textfield').fill(opts.password ?? devicePassword)
  await page.locator('#sign-in-button').click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
}

export async function logout(page: Page) {
  const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''
  const url = fullDomain ? `https://${fullDomain}/rest/logout` : '/rest/logout'
  await page.goto(url)
  await expect(page.locator('#username-textfield')).toBeVisible()
}
