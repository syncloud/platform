import { test, expect, Page } from '@playwright/test'
import { ssh } from '../helpers/ssh'
import { logout } from '../helpers/login'
import { uiMode } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
})

test.afterAll(async () => {
  await page.close()
})

test('permission denied for ldap user without access', async ({}, testInfo) => {
  const mode = uiMode(testInfo)
  ssh(`/snap/platform/current/openldap/bin/ldapadd.sh -x -w syncloud -D "dc=syncloud,dc=org" -f /test/test.${mode}.ldif`)
  await logout(page)
  await page.goto(`https://${fullDomain}`)
  await page.locator('#username-textfield').fill(`test${mode}`)
  await page.locator('#password-textfield').fill('password')
  await page.locator('#sign-in-button').click()
  const notif = page.locator('.notification')
  await expect(notif).toBeVisible()
  await expect(notif).toBeHidden({ timeout: 30_000 })
  await shoot(page, testInfo, 'permission-denied')
})
