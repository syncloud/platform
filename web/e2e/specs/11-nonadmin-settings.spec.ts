import { test, expect, Page } from '@playwright/test'
import { login, logout } from '../helpers/login'
import { ssh } from '../helpers/ssh'
import { settings } from '../helpers/ui'
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

test('regular user sees apps but no admin navigation', async ({}, testInfo) => {
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await expect(page.locator('#appcenter')).toHaveCount(0)
  await expect(page.locator('#appcenter_mobile')).toHaveCount(0)
  await shoot(page, testInfo, 'nav_nonadmin')
})

test('regular user manages only own passkeys', async ({}, testInfo) => {
  await settings(page, 'passkeys', testInfo)
  await expect(page.getByTestId('passkeys-title')).toBeVisible()
  await expect(page.getByTestId('passkeys-empty')).toBeVisible()
  await shoot(page, testInfo, 'passkeys_nonadmin')
})
