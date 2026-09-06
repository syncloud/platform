import { authenticator } from 'otplib'

const submitted = new Map<string, string>()

export function totp(secret: string): string {
  return authenticator.generate(secret)
}

export async function waitForFreshTotp(secret: string, minValidityMs = 5000): Promise<string> {
  const step = (authenticator.options.step ?? 30) * 1000
  for (;;) {
    const remaining = step - (Date.now() % step)
    const code = totp(secret)
    if (remaining >= minValidityMs && submitted.get(secret) !== code) {
      submitted.set(secret, code)
      return code
    }
    await new Promise(r => setTimeout(r, remaining + 1000))
  }
}
