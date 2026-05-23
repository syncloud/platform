import { Page, TestInfo } from '@playwright/test'
import * as path from 'node:path'
import * as fs from 'node:fs'

const artifactRoot = process.env.PLAYWRIGHT_ARTIFACT_DIR ?? 'artifact'

const freezeCss = '*,*::before,*::after { animation-duration: 0s !important; animation-delay: 0s !important; transition-duration: 0s !important; transition-delay: 0s !important; caret-color: transparent !important; }'

function viewOf(testInfo: TestInfo): string {
  return testInfo.project.name
}

async function waitForImages(page: Page) {
  await page.evaluate(() => Promise.all(
    Array.from(document.images).map(img =>
      img.complete && img.naturalWidth > 0
        ? Promise.resolve()
        : new Promise<void>(resolve => {
            img.addEventListener('load', () => resolve(), { once: true })
            img.addEventListener('error', () => resolve(), { once: true })
          })
    )
  ))
}

export async function shoot(page: Page, testInfo: TestInfo, name: string) {
  const view = viewOf(testInfo)
  const dir = path.join(artifactRoot, 'playwright', view, 'screenshot')
  fs.mkdirSync(dir, { recursive: true })
  const file = path.join(dir, `${name}-${view}.png`)
  await page.addStyleTag({ content: freezeCss })
  await waitForImages(page)
  await page.screenshot({ path: file, fullPage: false })
}
