import { test, expect, Page, TestInfo } from '@playwright/test'
import { login, deviceUser } from '../helpers/login'
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

async function openUsers (testInfo: TestInfo) {
  await settings(page, 'users', testInfo)
  await expect(page.getByTestId('users-title')).toBeVisible()
  await waitForLoading(page)
}

async function openUser (username: string) {
  await page.getByTestId('user-row-' + username).click()
  await expect(page.getByTestId('user-edit-title')).toBeVisible()
  await waitForLoading(page)
}

test('users list', async ({}, testInfo) => {
  await openUsers(testInfo)
  await expect(page.getByTestId('users-add')).toBeVisible()
  await defocus(page)
  await shoot(page, testInfo, 'settings_users')
})

test('last admin cannot be demoted', async ({}, testInfo) => {
  await openUsers(testInfo)
  await openUser(deviceUser)
  await expect(page.getByTestId('user-admin-last')).toBeVisible()
  await expect(page.locator('#user_admin')).toBeDisabled()
  await shoot(page, testInfo, 'settings_users_last_admin')
  await page.locator('#btn_cancel').click()
})

test('create requires username and a strong password', async ({}, testInfo) => {
  await openUsers(testInfo)
  await page.getByTestId('users-add').click()
  await expect(page.getByTestId('user-edit-title')).toBeVisible()
  await expect(page.locator('#btn_save')).toBeDisabled()

  await page.locator('#user_username').fill('temp')
  await expect(page.locator('#btn_save')).toBeDisabled()

  await page.locator('#user_password').fill('short')
  await expect(page.getByTestId('pwrule-length')).not.toHaveClass(/pw-ok/)
  await expect(page.getByTestId('pwrule-number')).not.toHaveClass(/pw-ok/)
  await expect(page.getByTestId('pwrule-letter')).toHaveClass(/pw-ok/)
  await expect(page.locator('#btn_save')).toBeDisabled()
  await shoot(page, testInfo, 'settings_users_password_rules')

  await page.locator('#user_password').fill('temppass1')
  await expect(page.getByTestId('pwrule-length')).toHaveClass(/pw-ok/)
  await expect(page.getByTestId('pwrule-letter')).toHaveClass(/pw-ok/)
  await expect(page.getByTestId('pwrule-number')).toHaveClass(/pw-ok/)
  await expect(page.locator('#btn_save')).toBeEnabled()

  await page.locator('#user_username').fill('')
  await expect(page.locator('#btn_save')).toBeDisabled()
  await page.locator('#btn_cancel').click()
  await waitForLoading(page)
})

test('add user with default email', async ({}, testInfo) => {
  await openUsers(testInfo)
  await page.getByTestId('users-add').click()
  await expect(page.getByTestId('user-edit-title')).toBeVisible()
  await page.locator('#user_username').fill('e2euser')
  await page.locator('#user_password').fill('e2epassword1')
  await shoot(page, testInfo, 'settings_users_create')
  await page.locator('#btn_save').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-row-e2euser')).toBeVisible()

  await openUser('e2euser')
  const email = await page.locator('#user_email').inputValue()
  expect(email.startsWith('e2euser@')).toBeTruthy()
  await page.locator('#btn_cancel').click()
})

test('change email and password', async ({}, testInfo) => {
  await openUsers(testInfo)
  await openUser('e2euser')
  await page.locator('#user_email').fill('e2euser@custom.org')
  await page.locator('#user_password').fill('changedpass1')
  await page.locator('#btn_save').click()
  await waitForLoading(page)

  await openUser('e2euser')
  await expect(page.locator('#user_email')).toHaveValue('e2euser@custom.org')
  await page.locator('#btn_cancel').click()
})

test('make user admin', async ({}, testInfo) => {
  await openUsers(testInfo)
  await openUser('e2euser')
  await expect(page.locator('#user_admin')).not.toBeChecked()
  await page.locator('xpath=//input[@id="user_admin"]/../span').click()
  await page.locator('#btn_save').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-admin-badge-e2euser')).toBeVisible()
  await defocus(page)
  await shoot(page, testInfo, 'settings_users_admin')
})

test('groups: syncloud hidden, create and assign inline', async ({}, testInfo) => {
  await openUsers(testInfo)
  await openUser('e2euser')
  await expect(page.getByTestId('user-group-syncloud')).toHaveCount(0)
  await page.locator('#new_group').fill('e2eteam')
  await page.getByTestId('group-create').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-group-e2eteam')).toHaveClass(/is-on/)
  await shoot(page, testInfo, 'settings_users_groups')
  await page.locator('#btn_save').click()
  await waitForLoading(page)

  await openUser('e2euser')
  await expect(page.getByTestId('user-group-e2eteam')).toHaveClass(/is-on/)
  await page.locator('#btn_cancel').click()
})

test('delete user asks for confirmation', async ({}, testInfo) => {
  await openUsers(testInfo)
  await openUser('e2euser')
  await page.locator('#btn_delete').click()
  await expect(page.getByTestId('btn_confirm')).toBeVisible()
  await shoot(page, testInfo, 'settings_users_delete_confirm')
  await page.getByTestId('btn_confirm').click()
  await waitForLoading(page)
  await expect(page.getByTestId('user-row-e2euser')).toHaveCount(0)
  await shoot(page, testInfo, 'settings_users_removed')
})
