import { test, expect, Page } from '@playwright/test'
import { login, logout } from '../helpers/login'
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

test('regular user sees apps but no admin navigation', async ({}, testInfo) => {
  await expect(page.locator('#apps')).toBeVisible()
  await expect(page.locator('#settings')).toHaveCount(0)
  await expect(page.locator('#appcenter')).toHaveCount(0)
  await shoot(page, testInfo, 'nav_nonadmin')
})
