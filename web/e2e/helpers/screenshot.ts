import { Page, TestInfo } from '@playwright/test'
import * as path from 'node:path'
import * as fs from 'node:fs'

const artifactRoot = process.env.PLAYWRIGHT_ARTIFACT_DIR ?? 'artifact'

function viewOf(testInfo: TestInfo): string {
  return testInfo.project.name
}

export async function shoot(page: Page, testInfo: TestInfo, name: string) {
  const view = viewOf(testInfo)
  const dir = path.join(artifactRoot, 'playwright', view, 'screenshot')
  fs.mkdirSync(dir, { recursive: true })
  const file = path.join(dir, `${name}-${view}.png`)
  await page.screenshot({ path: file, fullPage: false })
}
