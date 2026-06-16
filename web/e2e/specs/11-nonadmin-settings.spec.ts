import { test, expect, Page } from '@playwright/test'
import { login, logout } from '../helpers/login'
import { menu, settings } from '../helpers/ui'
import { ssh } from '../helpers/ssh'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const regularUser = 'regular'
const regularPassword = 'Regular123'

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  ssh(`snap run platform.cli user remove ${regularUser}`, { throw: false })
  ssh(`snap run platform.cli user add ${regularUser} --password=${regularPassword}`)
  await login(page, { user: regularUser, password: regularPassword })
})

test.afterAll(async () => {
  await logout(page)
  ssh(`snap run platform.cli user remove ${regularUser}`, { throw: false })
  await page.close()
})

test('settings hides admin tiles from a regular user', async ({}, testInfo) => {
  await menu(page, 'settings', testInfo)
  await expect(page.getByRole('heading', { name: 'Settings' })).toBeVisible()
  await expect(page.locator('#locale')).toBeVisible()
  await expect(page.locator('#twofactor')).toBeVisible()
  await expect(page.locator('#storage')).toHaveCount(0)
  await expect(page.locator('#users')).toHaveCount(0)
  await expect(page.locator('#customproxy')).toHaveCount(0)
  await expect(page.locator('#system')).toHaveCount(0)
  await shoot(page, testInfo, 'settings_nonadmin')
})

test('app center link hidden from a regular user', async ({}, testInfo) => {
  await expect(page.locator('#appcenter')).toHaveCount(0)
})

test('regular user can open locale', async ({}, testInfo) => {
  await settings(page, 'locale', testInfo)
  await expect(page.getByRole('heading', { name: 'Locale' })).toBeVisible()
  await shoot(page, testInfo, 'settings_locale_nonadmin')
})

test('regular user can open two-factor', async ({}, testInfo) => {
  await settings(page, 'twofactor', testInfo)
  await expect(page.getByRole('heading', { name: 'Two-Factor Authentication' })).toBeVisible()
  await shoot(page, testInfo, 'settings_twofactor_nonadmin')
})
