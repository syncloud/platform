import { test, expect, Page, CDPSession } from '@playwright/test'
import { env } from '../helpers/env'
import { logout, deviceUser, devicePassword } from '../helpers/login'
import { settings, waitForLoading, defocus } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

const fullDomain = env('PLAYWRIGHT_FULL_DOMAIN')

let page: Page
let cdp: CDPSession
let authenticatorId = ''

test.beforeAll(async ({ browser }) => {
  const context = await browser.newContext({ storageState: { cookies: [], origins: [] } })
  page = await context.newPage()
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
})

test.afterAll(async () => {
  await page.context().close()
})

test('passkey offered after password sign in', async ({}, testInfo) => {
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await page.getByTestId('username').fill(deviceUser)
  await page.getByTestId('password').fill(devicePassword)
  await page.route('**/api/secondfactor/webauthn/credential/register', async route => {
    const request = route.request()
    if (request.method() !== 'POST') {
      return route.continue()
    }
    const credential = request.postDataJSON()
    delete credential.clientExtensionResults.credProps
    await route.continue({ postData: JSON.stringify(credential) })
  })
  await page.getByTestId('sign-in').click()
  await expect(page.getByTestId('passkey-offer')).toBeVisible()
  await defocus(page)
  await shoot(page, testInfo, 'passkey-offer')
  await page.getByTestId('passkey-save').click()
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
  const credentials = await cdp.send('WebAuthn.getCredentials', { authenticatorId })
  expect(credentials.credentials.length).toBe(1)
  expect(credentials.credentials[0].isResidentCredential).toBe(true)
  await page.unroute('**/api/secondfactor/webauthn/credential/register')
})

test('passkey sign in from username autofill', async () => {
  await logout(page)
  const errorShown = page.getByTestId('error-dialog')
    .waitFor({ state: 'visible', timeout: 15_000 })
    .then(() => true)
    .catch(() => false)
  await page.goto(`https://auth.${fullDomain}`)
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
  expect(await errorShown).toBe(false)
})

test('passkey sign in button to OIDC app without autofill', async ({}, testInfo) => {
  await page.addInitScript(() => {
    PublicKeyCredential.isConditionalMediationAvailable = async () => false
  })
  await logout(page)
  await page.goto(`https://testapp.${fullDomain}/oidc/login`)
  await defocus(page)
  await shoot(page, testInfo, 'passkey-sign-in')
  await page.getByTestId('passkey-sign-in').click()
  await expect(page.locator('body')).toContainText(`OK ${deviceUser}`)
})

test('admin sees user passkeys', async ({}, testInfo) => {
  await page.goto('/')
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
  await settings(page, 'users', testInfo)
  await expect(page.getByTestId('users-title')).toBeVisible()
  await waitForLoading(page)
  await page.getByTestId('user-row-' + deviceUser).click()
  await expect(page.getByTestId('user-edit-title')).toBeVisible()
  await waitForLoading(page)
  await expect(page.getByTestId('user-passkeys').getByTestId(/^passkey-name-/)).toHaveCount(1)
  await shoot(page, testInfo, 'passkey-user-edit')
})

test('removed passkey no longer signs in', async ({}, testInfo) => {
  await settings(page, 'passkeys', testInfo)
  await expect(page.getByTestId('passkeys-title')).toBeVisible()
  await expect(page.getByTestId(/^passkey-name-/)).toHaveCount(1)
  await shoot(page, testInfo, 'passkeys')
  await page.getByTestId(/^passkey-remove-/).click()
  await page.getByTestId('btn_confirm').click()
  await expect(page.getByTestId('passkeys-empty')).toBeVisible()
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await page.getByTestId('passkey-sign-in').click()
  await expect(page.getByTestId('login-error')).toBeVisible()
})

test('password sign in offers a new passkey after a rejected one', async () => {
  await page.getByTestId('username').fill(deviceUser)
  await page.getByTestId('password').fill(devicePassword)
  await page.getByTestId('sign-in').click()
  await expect(page.getByTestId('passkey-offer')).toBeVisible()
  await page.getByTestId('passkey-skip').click()
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
})

test('skipped passkey offer is not repeated', async () => {
  await logout(page)
  await page.goto(`https://auth.${fullDomain}`)
  await page.getByTestId('username').fill(deviceUser)
  await page.getByTestId('password').fill(devicePassword)
  await page.getByTestId('sign-in').click()
  await expect(page.getByTestId('apps-title')).toBeVisible()
  await waitForLoading(page)
  await expect(page.getByTestId('passkey-offer')).toHaveCount(0)
})
