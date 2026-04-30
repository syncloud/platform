import { test, expect } from '@playwright/test'
import { login, waitForLoading } from '../helpers/login'
import { ssh } from '../helpers/ssh'
import { shoot } from '../helpers/screenshot'

const APP_ID = 'testapp'
const APP_SNAP = process.env.PLAYWRIGHT_TESTAPP_SNAP ?? '/test/testapp/testapp.snap'

test.describe('local app install', () => {
  test.beforeAll(async () => {
    if (ssh(`snap list ${APP_ID} 2>/dev/null`, { throw: false }).includes(APP_ID)) {
      ssh(`snap remove ${APP_ID}`)
    }
    ssh(`snap install --devmode --dangerous ${APP_SNAP}`)
  })

  test.afterAll(async () => {
    ssh(`snap remove ${APP_ID}`, { throw: false })
  })

  test('shows app page with local-install indicator and removes it', async ({ page }, testInfo) => {
    await login(page)
    await shoot(page, testInfo, 'local-app-applications')

    const tile = page.getByTestId(`app-tile-${APP_ID}`)
    await expect(tile).toBeVisible()
    await tile.click()

    await expect(page.getByTestId('app_name')).toBeVisible()
    await expect(page.getByTestId('local_install_badge')).toBeVisible()
    await expect(page.getByTestId('btn_remove')).toBeVisible()
    await waitForLoading(page)
    await shoot(page, testInfo, 'local-app-page')

    await page.getByTestId('btn_remove').click()
    await page.getByTestId('btn_confirm').click()
    await waitForLoading(page)

    await page.goto('/')
    await expect(page.getByTestId(`app-tile-${APP_ID}`)).toHaveCount(0)
    await shoot(page, testInfo, 'local-app-removed')
  })
})
