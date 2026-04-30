import { execFileSync } from 'node:child_process'

const host = process.env.PLAYWRIGHT_DEVICE_HOST ?? 'bookworm-amd64'
const sshKey = process.env.PLAYWRIGHT_SSH_KEY ?? '/root/.ssh/id_rsa'
const sshUser = process.env.PLAYWRIGHT_SSH_USER ?? 'root'

const baseArgs = [
  '-o', 'StrictHostKeyChecking=no',
  '-o', 'UserKnownHostsFile=/dev/null',
  '-o', 'LogLevel=ERROR',
  '-i', sshKey,
]

export function ssh(cmd: string, opts: { throw?: boolean } = {}): string {
  const args = [...baseArgs, `${sshUser}@${host}`, cmd]
  try {
    return execFileSync('ssh', args, { encoding: 'utf8', timeout: 120_000 })
  } catch (e: any) {
    if (opts.throw === false) {
      return e.stdout?.toString() ?? ''
    }
    throw e
  }
}
