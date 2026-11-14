import { test, expect, Page, CDPSession } from '@playwright/test'
import { login, logout, deviceUser, devicePassword } from '../helpers/login'
import { waitForLoading, defocus } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = process.env.PLAYWRIGHT_FULL_DOMAIN ?? process.env.PLAYWRIGHT_DOMAIN ?? ''

let page: Page
let cdp: CDPSession
let authenticatorId = ''

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  cdp = await page.context().newCDPSession(page)
  await cdp.send('WebAuthn.enable')
  const authenticator = await cdp.send('WebAuthn.addVirtualAuthenticator', {
    options: {
      protocol: 'ctap2',
      transport: 'internal',
      hasResidentKey: true,
      hasUserVerification: true,
      isUserVerified: true,
      automaticPresenceSimulation: true,
    },
  })
  authenticatorId = authenticator.authenticatorId
  await login(page)
})

test.afterAll(async () => {
  await page.close()
})

test('passkey setup on password sign in', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await page.getByTestId('username').fill(deviceUser)
  await page.getByTestId('password').fill(devicePassword)
  await page.getByTestId('passkey-setup').check()
  await defocus(page)
  await shoot(page, testInfo, 'passkey-setup')
  await page.getByTestId('sign-in').click()
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
  const credentials = await cdp.send('WebAuthn.getCredentials', { authenticatorId })
  expect(credentials.credentials.length).toBe(1)
  expect(credentials.credentials[0].isResidentCredential).toBe(true)
})

test('passkey sign in', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await defocus(page)
  await shoot(page, testInfo, 'passkey-sign-in')
  await page.getByTestId('passkey-sign-in').click()
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
})

test('passkey sign in to OIDC app', async () => {
  await logout(page)
  await page.goto(`https://testapp.${fullDomain}/oidc/login`)
  await page.getByTestId('passkey-sign-in').click()
  await expect(page.locator('body')).toContainText(`OK ${deviceUser}`)
})
