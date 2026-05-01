import { test, expect, Page } from '@playwright/test'
import { login } from '../helpers/login'
import { menu, settings, clickElSelect, waitForLoading } from '../helpers/ui'
import { ssh } from '../helpers/ssh'
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

test('settings root', async ({}, testInfo) => {
  await menu(page, 'settings', testInfo)
  await expect(page.getByRole('heading', { name: 'Settings' })).toBeVisible()
  await shoot(page, testInfo, 'settings')
})

test('settings activation', async ({}, testInfo) => {
  await settings(page, 'activation', testInfo)
  await expect(page.getByRole('heading', { name: 'Activation' })).toBeVisible()
  await shoot(page, testInfo, 'settings_activation')
})

test('settings access', async ({}, testInfo) => {
  await settings(page, 'access', testInfo)
  await expect(page.getByRole('heading', { name: 'Access' })).toBeVisible()
  await page.locator('xpath=//input[@id="tgl_ipv4_enabled"]/../span').click()
  await page.locator('css=#ipv4_mode_block[data-ready]').waitFor()
  await page.locator('xpath=//input[@id="tgl_ipv4_public"]/../span').click()
  await page.locator('css=#ipv4_public_block[data-ready]').waitFor()
  await shoot(page, testInfo, 'settings_access')
})

test('settings network', async ({}, testInfo) => {
  await settings(page, 'network', testInfo)
  await expect(page.getByRole('heading', { name: 'Network' })).toBeVisible()
  await shoot(page, testInfo, 'settings_network_unstable')
})

test('settings storage', async ({}, testInfo) => {
  await settings(page, 'storage', testInfo)
  await expect(page.getByRole('heading', { name: 'Storage' })).toBeVisible()
  await expect(page.locator('#btn_save')).toBeVisible()
  await shoot(page, testInfo, 'settings_storage_unstable')
})

test('settings updates', async ({}, testInfo) => {
  await settings(page, 'updates', testInfo)
  await expect(page.getByRole('heading', { name: 'Updates' })).toBeVisible()
  await shoot(page, testInfo, 'settings_updates_unstable')
})

test('settings internal memory', async ({}, testInfo) => {
  await settings(page, 'internalmemory', testInfo)
  await expect(page.getByRole('heading', { name: 'Internal Memory' })).toBeVisible()
  await shoot(page, testInfo, 'settings_internal_memory')
})

test('settings support', async ({}, testInfo) => {
  await settings(page, 'support', testInfo)
  await expect(page.getByRole('heading', { name: 'Support' })).toBeVisible()
  await shoot(page, testInfo, 'settings_support')
})

test('settings backup', async ({}, testInfo) => {
  await settings(page, 'backup', testInfo)
  await expect(page.getByRole('heading', { name: 'Backup' })).toBeVisible()
  await shoot(page, testInfo, 'settings_backup_unstable')
  await expect(page.locator('.el-notification__title')).toHaveCount(0)
  await clickElSelect(page, 'auto')
  await page.locator('#auto-backup').click()
  await clickElSelect(page, 'auto-day')
  await page.locator('#auto-day-monday').click()
  await clickElSelect(page, 'auto-hour')
  await page.locator('#auto-hour-1').click()
  await page.locator('#save').click()
  await shoot(page, testInfo, 'settings_backup_saved_unstable')
  await expect(page.locator('.el-notification__title')).toHaveCount(0)
})

test('settings certificate', async ({}, testInfo) => {
  await settings(page, 'certificate', testInfo)
  await expect(page.getByRole('heading', { name: 'Certificate' })).toBeVisible()
  await shoot(page, testInfo, 'settings_certificate')
})

test('settings locale', async ({}, testInfo) => {
  await settings(page, 'locale', testInfo)
  await expect(page.getByRole('heading', { name: 'Locale' })).toBeVisible()
  await shoot(page, testInfo, 'settings_locale_unstable')
  const baseline = await page.locator('#current_time').textContent()
  await clickElSelect(page, 'settings_timezone')
  await page.keyboard.type('Asia/Tokyo')
  await page.waitForTimeout(1000)
  await page.locator('#settings_tz_Asia_Tokyo').click()
  await page.locator('#btn_save_timezone').click()
  await waitForLoading(page)
  const tz = ssh('cat /etc/timezone').trim()
  expect(tz).toBe('Asia/Tokyo')
  await page.waitForTimeout(2000)
  const updated = await page.locator('#current_time').textContent()
  expect(updated).not.toBe(baseline)
  await shoot(page, testInfo, 'settings_locale_tokyo_unstable')
})
