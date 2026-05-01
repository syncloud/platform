import { authenticator } from 'otplib'

export function totp(secret: string): string {
  return authenticator.generate(secret)
}

export async function waitForFreshTotp(secret: string): Promise<string> {
  const step = (authenticator.options.step ?? 30) * 1000
  const remaining = step - (Date.now() % step)
  await new Promise(r => setTimeout(r, remaining + 1000))
  return totp(secret)
}
