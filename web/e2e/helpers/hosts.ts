import { execFileSync, execSync } from 'node:child_process'
import * as dns from 'node:dns/promises'

export async function addHostAlias(alias: string, deviceHost: string, baseDomain: string): Promise<void> {
  const fqdn = `${alias}.${baseDomain}`
  const records = await dns.lookup(deviceHost).catch(() => null)
  if (!records) throw new Error(`could not resolve ${deviceHost}`)
  const line = `${records.address} ${fqdn}`
  try {
    execSync(`grep -qF "${line}" /etc/hosts || echo "${line}" | sudo tee -a /etc/hosts > /dev/null`, { stdio: 'inherit' })
  } catch {
    execFileSync('sh', ['-c', `echo "${line}" >> /etc/hosts`])
  }
}
