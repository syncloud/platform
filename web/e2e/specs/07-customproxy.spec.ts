import { test, expect, request, Page } from '@playwright/test'
import { addHostAlias } from '../helpers/hosts'
import { ssh, deviceHost } from '../helpers/ssh'
import { login } from '../helpers/login'
import { settings, waitForLoading, defocus } from '../helpers/ui'
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

test('custom proxy overrides missing app (files alias)', async ({}, testInfo) => {
  ssh('nohup /test/externalapp/externalapp > /tmp/syncloud/ui/externalapp.log 2>&1 &', { throw: false })
  await addHostAlias('files', deviceHost, fullDomain)
  await settings(page, 'customproxy', testInfo)
  await expect(page.getByRole('heading', { name: 'Custom Proxy' })).toBeVisible()
  await page.locator('#proxy_name').fill('files')
  await page.locator('#proxy_host').fill('localhost')
  await page.locator('#proxy_port').fill('8585')
  await page.locator('#btn_add').click()
  await waitForLoading(page)
  await defocus(page)
  await shoot(page, testInfo, 'settings_custom_proxy_files_added')

  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  await expect.poll(async () => {
    const r = await ctx.get(`https://files.${fullDomain}`)
    if (r.status() !== 200) return 'pending'
    return await r.text()
  }, { timeout: 60_000 }).toBe('external')
  await ctx.dispose()
  await shoot(page, testInfo, 'settings_custom_proxy_files_verified')
})

test('remove custom proxy files alias', async ({}, testInfo) => {
  await settings(page, 'customproxy', testInfo)
  await expect(page.getByRole('heading', { name: 'Custom Proxy' })).toBeVisible()
  await page.locator('#btn_remove_files').click()
  await waitForLoading(page)
  await shoot(page, testInfo, 'settings_custom_proxy_files_removed')
})

test('custom proxy with Authelia redirects to login', async ({}, testInfo) => {
  ssh('nohup /test/externalapp/externalapp 8586 > /tmp/syncloud/ui/externalapp-protected.log 2>&1 &', { throw: false })
  await addHostAlias('protected', deviceHost, fullDomain)
  await settings(page, 'customproxy', testInfo)
  await expect(page.getByRole('heading', { name: 'Custom Proxy' })).toBeVisible()
  await page.locator('#proxy_name').fill('protected')
  await page.locator('#proxy_host').fill('localhost')
  await page.locator('#proxy_port').fill('8586')
  await page.getByTestId('proxy-authelia').check()
  await shoot(page, testInfo, 'settings_custom_proxy_authelia_filled')
  await page.locator('#btn_add').click()
  await waitForLoading(page)
  await defocus(page)
  await shoot(page, testInfo, 'settings_custom_proxy_authelia_added')

  await expect(page.getByTestId('proxy-row-protected-authelia')).toBeVisible()

  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  await expect.poll(async () => {
    const r = await ctx.get(`https://protected.${fullDomain}`, { maxRedirects: 0 })
    return r.status()
  }, { timeout: 60_000 }).toBe(302)
  const r = await ctx.get(`https://protected.${fullDomain}`, { maxRedirects: 0 })
  expect(r.headers()['location']).toContain(`auth.${fullDomain}`)
  await ctx.dispose()

  await settings(page, 'customproxy', testInfo)
  await page.locator('#btn_remove_protected').click()
  await waitForLoading(page)
})

test('custom proxy externalapp', async ({}, testInfo) => {
  await addHostAlias('externalapp', deviceHost, fullDomain)
  await settings(page, 'customproxy', testInfo)
  await expect(page.getByRole('heading', { name: 'Custom Proxy' })).toBeVisible()
  await waitForLoading(page)
  await shoot(page, testInfo, 'settings_custom_proxy')
  await page.locator('#proxy_name').fill('externalapp')
  await page.locator('#proxy_host').fill('localhost')
  await page.locator('#proxy_port').fill('8585')
  await shoot(page, testInfo, 'settings_custom_proxy_filled')
  await page.locator('#btn_add').click()
  await waitForLoading(page)
  await expect(page.locator('a', { hasText: 'externalapp' })).toBeVisible()
  await defocus(page)
  await shoot(page, testInfo, 'settings_custom_proxy_added')

  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  await expect.poll(async () => {
    const r = await ctx.get(`https://externalapp.${fullDomain}`)
    if (r.status() !== 200) return 'pending'
    return await r.text()
  }, { timeout: 60_000 }).toBe('external')
  await ctx.dispose()

  await shoot(page, testInfo, 'settings_custom_proxy_verified')
})
