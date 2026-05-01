import { request, APIRequestContext } from '@playwright/test'
import { deviceHost, ssh } from './ssh'

const deviceUser = process.env.PLAYWRIGHT_DEVICE_USER ?? 'user'
const devicePassword = process.env.PLAYWRIGHT_DEVICE_PASSWORD ?? 'Password1'

export async function loginV2(): Promise<APIRequestContext> {
  const token = ssh(`snap run platform.cli login ${deviceUser} ${devicePassword}`).trim()
  const ctx = await request.newContext({
    baseURL: `https://${deviceHost}`,
    ignoreHTTPSErrors: true,
  })
  const resp = await ctx.post('/rest/login/token', { data: { token } })
  if (!resp.ok()) {
    await ctx.dispose()
    throw new Error(`token login failed: ${resp.status()} ${await resp.text()}`)
  }
  return ctx
}

export async function activated(): Promise<boolean> {
  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  try {
    const resp = await ctx.get(`https://${deviceHost}/rest/activation/status`)
    if (!resp.ok()) return false
    const body = await resp.json()
    return body?.data?.activated === true
  } catch {
    return false
  } finally {
    await ctx.dispose()
  }
}

export async function waitForRest(url: string, status = 200, timeoutSec = 60): Promise<void> {
  const ctx = await request.newContext({ ignoreHTTPSErrors: true })
  const deadline = Date.now() + timeoutSec * 1000
  try {
    while (Date.now() < deadline) {
      try {
        const r = await ctx.get(url)
        if (r.status() === status) return
      } catch {}
      await new Promise(r => setTimeout(r, 1000))
    }
    throw new Error(`waitForRest timeout ${url}`)
  } finally {
    await ctx.dispose()
  }
}
