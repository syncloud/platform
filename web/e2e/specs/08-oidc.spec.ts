import { test, expect, Page } from '@playwright/test'
import { login, logout, deviceUser, devicePassword } from '../helpers/login'
import { waitForLoading, defocus } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  await login(page)
})

test.afterAll(async () => {
  await page.close()
})

test('testapp OIDC login', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://testapp.${fullDomain}/oidc/login`)
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await defocus(page)
  await shoot(page, testInfo, 'testapp-oidc-login')
  await page.locator('#sign-in-button').click()
  await expect(page.locator(`xpath=//*[contains(text(),'OK ${deviceUser}')]`)).toBeVisible()
  await shoot(page, testInfo, 'testapp-oidc-callback')
})

test('auth web (auth.domain → main app)', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await defocus(page)
  await shoot(page, testInfo, 'auth')
  await page.locator('#sign-in-button').click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
})
