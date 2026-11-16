import { defineConfig, devices } from '@playwright/test'
import { env } from './helpers/env'

const domain = env('PLAYWRIGHT_DOMAIN')
const app = env('PLAYWRIGHT_APP')
const artifactDir = env('PLAYWRIGHT_ARTIFACT_DIR')
const fullDomain = env('PLAYWRIGHT_FULL_DOMAIN')

export default defineConfig({
  testDir: './specs',
  fullyParallel: false,
  workers: 1,
  retries: 0,
  maxFailures: 1,
  reporter: [['list'], ['html', { open: 'never', outputFolder: `${artifactDir}/playwright/report` }]],
  outputDir: `${artifactDir}/playwright/test-results`,
  globalTeardown: './globalTeardown.ts',
  timeout: 60_000,
  expect: { timeout: 10_000 },
  use: {
    baseURL: `https://${domain}`,
    storageState: {
      cookies: [],
      origins: [{
        origin: `https://auth.${fullDomain}`,
        localStorage: [{ name: 'passkey', value: `declined:${Date.now()}` }],
      }],
    },
    ignoreHTTPSErrors: true,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    launchOptions: {
      args: [
        '--font-render-hinting=none',
        '--disable-skia-runtime-opts',
        '--force-color-profile=srgb',
      ],
    },
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
