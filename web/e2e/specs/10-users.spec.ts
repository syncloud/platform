import { test, expect, Page } from '@playwright/test'
import { login } from '../helpers/login'
import { settings, waitForLoading, defocus } from '../helpers/ui'
import { shoot } from '../helpers/screenshot'

test.describe.configure({ mode: 'serial' })

let page: Page

test.beforeAll(async ({ browser }) => {
  page = await browser.newPage()
  await login(page)
})

test.afterAll(async () => {
  await page.close()
})

test('users page', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await expect(page.getByTestId('users-title')).toBeVisible()
  await waitForLoading(page)
  await defocus(page)
  await shoot(page, testInfo, 'settings_users')
})

test('add user with default email', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await page.locator('#user_username').fill('testuser')
  await page.locator('#user_password').fill('testpassword')
  await shoot(page, testInfo, 'settings_users_filled')
  await page.locator('#btn_add_user').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-row-testuser')).toBeVisible()

  const email = await page.getByTestId('user-email-testuser').inputValue()
  expect(email.startsWith('testuser@')).toBeTruthy()
  await defocus(page)
  await shoot(page, testInfo, 'settings_users_added')
})

test('change user email', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await page.getByTestId('user-email-testuser').fill('testuser@custom.org')
  await page.locator('#btn_email_testuser').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-email-testuser')).toHaveValue('testuser@custom.org')
})

test('make user admin', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await expect(page.locator('#user_admin_testuser')).not.toBeChecked()
  await page.locator('xpath=//input[@id="user_admin_testuser"]/../span').click()
  await waitForLoading(page)
  await expect(page.locator('#user_admin_testuser')).toBeChecked()
  await defocus(page)
  await shoot(page, testInfo, 'settings_users_admin')
})

test('add group and membership', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await page.locator('#group_name').fill('team')
  await page.locator('#btn_add_group').click()
  await waitForLoading(page)
  await expect(page.getByTestId('group-row-team')).toBeVisible()

  await page.getByTestId('user-group-testuser-team').check()
  await waitForLoading(page)
  await expect(page.getByTestId('user-group-testuser-team')).toBeChecked()
  await defocus(page)
  await shoot(page, testInfo, 'settings_users_group')
})

test('remove group and user', async ({}, testInfo) => {
  await settings(page, 'users', testInfo)
  await page.locator('#btn_remove_group_team').click()
  await waitForLoading(page)
  await expect(page.getByTestId('group-row-team')).toHaveCount(0)

  await page.locator('#btn_remove_user_testuser').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-row-testuser')).toHaveCount(0)
  await shoot(page, testInfo, 'settings_users_removed')
})
