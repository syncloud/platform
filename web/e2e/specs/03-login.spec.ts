import { test, expect, Page } from '@playwright/test'
import { login, deviceUser, devicePassword, logout } from '../helpers/login'
import { waitForLoading, waitAppIconsLoaded, defocus } from '../helpers/ui'
import { ssh } from '../helpers/ssh'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
})

test.afterAll(async () => {
  await page.close()
})

test('login page loads on auth subdomain', async ({}, testInfo) => {
  const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''
  await page.goto(`https://auth.${fullDomain}`)
  await expect(page.locator('#username-textfield')).toBeVisible()
  await shoot(page, testInfo, 'login-page-direct')
})

test('login', async ({}, testInfo) => {
  await page.goto('/')
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await defocus(page)
  await shoot(page, testInfo, 'login')
  await page.locator('#sign-in-button').click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
  await expect(page.locator('.appimg').first()).toBeVisible()
  await waitAppIconsLoaded(page)
  await shoot(page, testInfo, 'index')
})

test('regular user login (no 2FA enabled)', async ({}, testInfo) => {
  ssh('snap run platform.cli user add regularuser --password=regularpass123', { throw: false })
  try {
    await logout(page)
    await page.goto('/')
    await page.locator('#username-textfield').fill('regularuser')
    await page.locator('#password-textfield').fill('regularpass123')
    await page.locator('#sign-in-button').click()
    await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
    await waitForLoading(page)
    await shoot(page, testInfo, 'login-regular-user')
  } finally {
    ssh('snap run platform.cli user remove regularuser', { throw: false })
    await logout(page)
    await login(page)
  }
})
