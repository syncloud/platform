import { test, expect, Page } from '@playwright/test'
import { ssh, deviceHost } from '../helpers/ssh'
import { settings, waitForLoading, waitAppIconsLoaded, defocus } from '../helpers/ui'
import { login, deviceUser, devicePassword } from '../helpers/login'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''
const domain = process.env.PLAYWRIGHT_DOMAIN_SHORT ?? deviceHost.replace(/\.redirect$/, '')
const redirectUser = process.env.PLAYWRIGHT_REDIRECT_USER ?? 'redirect'
const redirectPassword = process.env.PLAYWRIGHT_REDIRECT_PASSWORD ?? 'redirect'

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  await login(page)
})

test.afterAll(async () => {
  await page.close()
})

test('deactivate via settings then re-activate', async ({}, testInfo) => {
  await page.goto(`https://${fullDomain}`)
  await settings(page, 'activation', testInfo)
  await expect(page.getByRole('heading', { name: 'Activation' })).toBeVisible()
  await page.locator('#btn_reactivate').click()
  await page.locator('#btn_welcome_next').click()
  await expect(page.locator('#btn_syncloud_domain')).toBeVisible()
  await shoot(page, testInfo, 'activate-empty')
  await page.locator('#btn_syncloud_domain').click()
  await expect(page.locator('#email')).toBeVisible()
  await page.locator('#email').fill(redirectUser)
  await defocus(page)
  await shoot(page, testInfo, 'activate-redirect-email')
  await page.locator('#redirect_password').fill(redirectPassword)
  await defocus(page)
  await shoot(page, testInfo, 'activate-account')
  await page.locator('#btn_account_next').click()
  await expect(page.locator('#domain_input')).toBeVisible()
  await page.locator('#domain_input').fill(domain)
  await defocus(page)
  await shoot(page, testInfo, 'activate-type')
  await page.locator('#btn_domain_next').click()
  await waitForLoading(page)
  await expect(page.locator('#btn_activate')).toBeVisible()
  await shoot(page, testInfo, 'activate-redirect')
  await page.locator('#device_username').fill(deviceUser)
  await page.locator('#device_password').fill(devicePassword)
  await page.locator('#device_password_confirm').fill(devicePassword)
  await defocus(page)
  await shoot(page, testInfo, 'activate-ready')
  await page.locator('#btn_activate').click()
  await waitForLoading(page)
  await expect(page.locator('#username-textfield')).toBeVisible()
  await defocus(page)
  await shoot(page, testInfo, 'deactivate-login-page')
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await page.locator('#sign-in-button').click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
  await expect(page.locator('.appimg').first()).toBeVisible()
  await waitAppIconsLoaded(page)
  await shoot(page, testInfo, 'reactivate-index')
})
