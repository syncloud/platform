import { test, expect, Page } from '@playwright/test'
import { login, logout, deviceUser, devicePassword } from '../helpers/login'
import { ssh } from '../helpers/ssh'
import { settings, waitForLoading, waitAppIconsLoaded } from '../helpers/ui'
import { loginV2 } from '../helpers/device'
import { waitForFreshTotp } from '../helpers/totp'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''

let page: Page
let storedTotpSecret = ''

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  await login(page)
})

test.afterAll(async () => {
  await page.close()
})

test('2FA settings page', async ({}, testInfo) => {
  await page.goto(`https://${fullDomain}`)
  await settings(page, 'twofactor', testInfo)
  await expect(page.getByRole('heading', { name: 'Two-Factor Authentication' })).toBeVisible()
  await expect(page.locator('#twofa_status')).toBeVisible()
  await shoot(page, testInfo, '2fa_settings')
})

test('2FA enable', async ({}, testInfo) => {
  await page.goto(`https://${fullDomain}`)
  await settings(page, 'twofactor', testInfo)
  await page.locator('#btn_enable_2fa').click()
  await expect(page.locator('#btn_disable_2fa')).toBeVisible({ timeout: 60_000 })
  await shoot(page, testInfo, '2fa_enabled_unstable')
})

test('2FA login (first time, captures TOTP secret)', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://${fullDomain}`)
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await page.locator('#sign-in-button').click()
  await expect(page.locator('#totp_qr')).toBeVisible()
  await shoot(page, testInfo, '2fa_enabled_first_login_unstable')
  storedTotpSecret = (await page.locator('#totp_secret').textContent())?.trim() ?? ''
  expect(storedTotpSecret.length).toBeGreaterThan(0)

  const code = await waitForFreshTotp(storedTotpSecret)
  await page.locator('#otp-input').fill(code)
  await shoot(page, testInfo, '2fa_login_totp_unstable')
  await page.locator("button[type='button']").click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
  await expect(page.locator('.appimg').first()).toBeVisible()
  await waitAppIconsLoaded(page)
  await shoot(page, testInfo, '2fa_login_success')
})

test('2FA login returning user (TOTP only)', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://${fullDomain}`)
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await page.locator('#sign-in-button').click()
  await expect(page.locator('#otp-input')).toBeVisible()
  await shoot(page, testInfo, '2fa_returning')
  const code = await waitForFreshTotp(storedTotpSecret)
  await page.locator('#otp-input').fill(code)
  await page.locator("button[type='button']").click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
  await waitAppIconsLoaded(page)
  await shoot(page, testInfo, '2fa_login_success')
})

test('2FA regular user login (own QR)', async ({}, testInfo) => {
  ssh('snap run platform.cli user add testuser --password=testpass123', { throw: false })
  try {
    await logout(page)
    await page.goto(`https://${fullDomain}`)
    await page.locator('#username-textfield').fill('testuser')
    await page.locator('#password-textfield').fill('testpass123')
    await page.locator('#sign-in-button').click()

    await expect(page.locator('#totp_qr')).toBeVisible()
    const userSecret = (await page.locator('#totp_secret').textContent())?.trim() ?? ''
    const code = await waitForFreshTotp(userSecret)
    await page.locator('#otp-input').fill(code)
    await page.locator("button[type='button']").click()

    await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
    await waitForLoading(page)
    await waitAppIconsLoaded(page)
    await shoot(page, testInfo, '2fa_regular_user_login')
  } finally {
    ssh('snap run platform.cli user remove testuser', { throw: false })
    await logout(page)
    await page.goto(`https://${fullDomain}`)
    await page.locator('#username-textfield').fill(deviceUser)
    await page.locator('#password-textfield').fill(devicePassword)
    await page.locator('#sign-in-button').click()
    await expect(page.locator('#otp-input')).toBeVisible()
    const code = await waitForFreshTotp(storedTotpSecret)
    await page.locator('#otp-input').fill(code)
    await page.locator("button[type='button']").click()
    await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
    await waitForLoading(page)
  }
})

test('2FA disable', async ({}, testInfo) => {
  await settings(page, 'twofactor', testInfo)
  await expect(page.getByRole('heading', { name: 'Two-Factor Authentication' })).toBeVisible()
  await page.locator('#btn_disable_2fa').click()
  await shoot(page, testInfo, '2fa_disabled')
})

test('2FA recovery via CLI', async ({}, testInfo) => {
  const ctx = await loginV2()
  await ctx.post(`https://${fullDomain}/rest/settings/2fa`, {
    data: { enabled: true },
  })
  await ctx.dispose()

  ssh('snap run platform.cli disable-2fa')
  await page.waitForTimeout(2000)

  await logout(page)
  await page.goto(`https://${fullDomain}`)
  await page.locator('#username-textfield').fill(deviceUser)
  await page.locator('#password-textfield').fill(devicePassword)
  await page.locator('#sign-in-button').click()
  await expect(page.getByRole('heading', { name: 'Applications' })).toBeVisible()
  await waitForLoading(page)
  await waitAppIconsLoaded(page)
  await shoot(page, testInfo, '2fa_recovery_cli')
})
