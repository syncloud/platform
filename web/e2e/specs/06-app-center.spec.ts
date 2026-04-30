import { test, expect, request, Page } from '@playwright/test'
import { login } from '../helpers/login'
import { menu, waitForLoading } from '../helpers/ui'
import { ssh } from '../helpers/ssh'
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

test('app center', async ({}, testInfo) => {
  await menu(page, 'appcenter', testInfo)
  await expect(page.getByRole('heading', { name: 'App Center' })).toBeVisible()
  await expect(page.getByText('File browser')).toBeVisible()
  await page.locator('#appcenter_filter').fill('nextcloud')
  await expect(page.getByText('Nextcloud file sharing')).toBeVisible()
  await shoot(page, testInfo, 'appcenter')
})

test('install app', async ({}, testInfo) => {
  await menu(page, 'appcenter', testInfo)
  await expect(page.getByRole('heading', { name: 'App Center' })).toBeVisible()
  await page.getByText('File browser').click()
  await expect(page.getByRole('heading', { name: 'File browser' })).toBeVisible()
  await page.locator('#btn_install').click()
  await page.locator('#btn_confirm').click()
  await expect(page.locator('#btn_remove')).toBeVisible({ timeout: 600_000 })
  await shoot(page, testInfo, 'app_installed')

  ssh('ls -la /var/snap/files/common/', { throw: false })

  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  await expect.poll(async () => {
    const r = await ctx.get(`https://files.${fullDomain}`)
    if (r.status() !== 200) return 'pending'
    const text = await r.text()
    return text === 'external' ? 'external' : 'real'
  }, { timeout: 60_000, intervals: [1000, 2000, 5000] }).toBe('real')
  await ctx.dispose()
})

test('remove app', async ({}, testInfo) => {
  await page.locator('#btn_remove').click()
  await page.locator('#btn_confirm').click()
  await waitForLoading(page)
  await expect(page.locator('#btn_install')).toBeVisible()
  await shoot(page, testInfo, 'app_removed')
})

test('not installed app', async ({}, testInfo) => {
  await menu(page, 'appcenter', testInfo)
  await page.getByText('Nextcloud file sharing').click()
  await expect(page.getByRole('heading', { name: 'Nextcloud file sharing' })).toBeVisible()
  await expect(page.locator('#btn_install')).toBeVisible()
  await shoot(page, testInfo, 'app_not_installed')
})

test('502 page', async ({}, testInfo) => {
  await page.goto(`https://unknown.${fullDomain}`)
  await expect(page.getByRole('heading', { name: /App is not available/ })).toBeVisible()
})
