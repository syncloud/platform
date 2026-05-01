import { test, expect, request, Page } from '@playwright/test'
import { addHostAlias } from '../helpers/hosts'
import { ssh, deviceHost } from '../helpers/ssh'
import { loginV2, waitForRest } from '../helpers/device'
import { waitForLoading, defocus } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'
import { deviceUser, devicePassword } from '../helpers/login'

test.describe.configure({ mode: 'serial' })

const domain = process.env.PLAYWRIGHT_DOMAIN_SHORT ?? deriveShortDomain()
const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? `${domain}.redirect`
const mainDomain = process.env.PLAYWRIGHT_MAIN_DOMAIN ?? 'redirect'
const redirectUser = process.env.PLAYWRIGHT_REDIRECT_USER ?? 'redirect'
const redirectPassword = process.env.PLAYWRIGHT_REDIRECT_PASSWORD ?? 'redirect'

function deriveShortDomain(): string {
  const d = process.env.PLAYWRIGHT_DOMAIN ?? deviceHost
  return d.replace(/\.redirect$/, '')
}

const LOCALES = ['en', 'zh-CN', 'es', 'hi', 'ar', 'pt', 'ru', 'ja', 'de', 'fr']

let page: Page

test.beforeAll(async ({ browser }) => {
  await addHostAlias(process.env.PLAYWRIGHT_APP ?? 'platform', deviceHost, domain)
  await addHostAlias('auth', deviceHost, fullDomain)
  await addHostAlias('testapp', deviceHost, fullDomain)
  page = await browser.newPage()
})

test.afterAll(async () => {
  await page.close()
})

test('deactivate', async () => {
  ssh(`echo "$(getent hosts ${deviceHost} | awk '{print $1}') auth.${fullDomain}" >> /etc/hosts`, { throw: false })
  ssh('snap run platform.cli config set redirect.domain ' + mainDomain)
  ssh('snap run platform.cli config set certbot.staging true')
  ssh('snap run platform.cli config set redirect.api_url http://api.redirect')

  const ctx = await loginV2()
  const resp = await ctx.post(`https://${deviceHost}/rest/deactivate`)
  expect(resp.status()).toBe(200)
  expect(await resp.text()).toContain('"success":true')
  await ctx.dispose()
})

test('fake cert', async ({}, testInfo) => {
  ssh('rm /var/snap/platform/current/syncloud.crt', { throw: false })
  ssh('snap run platform.cli cert')
  ssh('snap restart platform')
  await waitForRest(`https://${deviceHost}/rest/activation/status`, 200, 60)
  await page.goto(`https://${deviceHost}`)
  await expect(page.locator('#btn_welcome_next')).toBeVisible()
  await waitForLoading(page)
  await shoot(page, testInfo, 'fake-cert_unstable')
})

test('login page loads on auth subdomain', async ({}, testInfo) => {
  await page.goto(`https://auth.${fullDomain}`)
  await expect(page.locator('#username-textfield')).toBeVisible()
  await shoot(page, testInfo, 'login-page-direct')
})

test('activate languages', async ({}, testInfo) => {
  await page.goto(`https://${deviceHost}`)
  await expect(page.locator('h1').first()).toBeVisible()
  await waitForLoading(page)
  for (const code of LOCALES) {
    await page.evaluate((c) => {
      window.localStorage.setItem('syncloud.locale', c)
      window.location.reload()
    }, code)
    await expect(page.locator('h1').first()).toBeVisible()
    await waitForLoading(page)
    await shoot(page, testInfo, `activate-lang-${code}_unstable`)
  }
  await page.evaluate(() => {
    window.localStorage.setItem('syncloud.locale', 'en')
    window.location.reload()
  })
  await expect(page.locator('#btn_welcome_next')).toBeVisible()
  await waitForLoading(page)
})

test('activate', async ({}, testInfo) => {
  test.setTimeout(120_000)
  await page.goto(`https://${deviceHost}`)
  await expect(page.locator('#btn_welcome_next')).toBeVisible()
  await shoot(page, testInfo, 'activate-welcome_unstable')
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
  await expect(page.locator('#device_username')).toBeVisible()
  await expect(page.locator('#device_password')).toBeVisible()
  await page.locator('#device_username').fill(deviceUser)
  await page.locator('#device_password').fill(devicePassword)
  await page.locator('#device_password_confirm').fill(devicePassword)
  await defocus(page)
  await shoot(page, testInfo, 'activate-ready')
  await page.locator('#btn_activate').click()
  await waitForLoading(page)
  await expect(page.locator('#username-textfield')).toBeVisible({ timeout: 30_000 })
  await defocus(page)
  await shoot(page, testInfo, 'activate')
})
