import { Page, TestInfo, expect } from '@playwright/test'

export type UiMode = 'desktop' | 'mobile'

export function uiMode(testInfo: TestInfo): UiMode {
  return testInfo.project.name as UiMode
}

export async function waitForLoading(page: Page) {
  await page.waitForTimeout(500)
  const mask = page.locator('.el-loading-mask').first()
  await mask.waitFor({ state: 'hidden', timeout: 120_000 }).catch(() => {})
}

export async function defocus(page: Page) {
  await page.locator('body').click({ position: { x: 1, y: 1 } })
  await page.waitForTimeout(500)
}

export async function menu(page: Page, id: string, testInfo: TestInfo) {
  const mode = uiMode(testInfo)
  const target = mode === 'mobile' ? `${id}_mobile` : id
  for (let attempt = 0; attempt < 10; attempt++) {
    try {
      if (mode === 'mobile') {
        await page.locator('#menubutton').click()
        await expect(page.locator(`#${target}`)).toBeVisible()
      }
      await page.locator(`#${target}`).click()
      if (mode === 'mobile') {
        await expect(page.locator('#menu')).toBeHidden({ timeout: 10_000 })
      }
      await waitForLoading(page)
      return
    } catch (e) {
      if (attempt === 9) throw e
      await page.waitForTimeout(1000)
    }
  }
}

export async function settings(page: Page, key: string, testInfo: TestInfo) {
  await menu(page, 'settings', testInfo)
  await page.locator(`#${key}`).click()
  await waitForLoading(page)
}

export async function clickElSelect(page: Page, selectId: string) {
  const wrapper = page.locator(`xpath=//*[@id='${selectId}']/ancestor::div[contains(@class, 'el-select__wrapper')]`)
  await wrapper.click()
}

export async function waitAppIconsLoaded(page: Page, timeoutMs = 10_000) {
  await page.waitForFunction(() => {
    const imgs = Array.from(document.querySelectorAll<HTMLImageElement>('.appimg'))
    return imgs.length > 0 && imgs.every(i => i.complete && i.naturalWidth > 0)
  }, undefined, { timeout: timeoutMs })
}
