import { defineConfig, devices } from '@playwright/test'

const domain = process.env.PLAYWRIGHT_DOMAIN ?? 'bookworm-amd64.redirect'
const app = process.env.PLAYWRIGHT_APP ?? 'platform'
const artifactDir = process.env.PLAYWRIGHT_ARTIFACT_DIR ?? 'artifact'

export default defineConfig({
  testDir: './specs',
  fullyParallel: false,
  workers: 1,
  retries: 0,
  maxFailures: 1,
  reporter: [['list'], ['html', { open: 'never', outputFolder: `${artifactDir}/playwright-report` }]],
  outputDir: `${artifactDir}/test-results`,
  globalTeardown: './globalTeardown.ts',
  timeout: 60_000,
  expect: { timeout: 10_000 },
  use: {
    baseURL: `https://${domain}`,
    ignoreHTTPSErrors: true,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  projects: [
    {
      name: 'desktop',
      use: { ...devices['Desktop Chrome'], viewport: { width: 1440, height: 960 } },
    },
    {
      name: 'mobile',
      use: { ...devices['Desktop Chrome'], viewport: { width: 390, height: 844 } },
    },
  ],
  metadata: {
    app,
    domain,
    artifactDir,
  },
})
