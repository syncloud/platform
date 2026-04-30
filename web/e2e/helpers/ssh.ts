import { execFileSync } from 'node:child_process'

export const deviceHost = process.env.PLAYWRIGHT_DEVICE_HOST ?? 'bookworm-amd64'
export const sshKey = process.env.PLAYWRIGHT_SSH_KEY ?? '/root/.ssh/id_rsa'
export const sshUser = process.env.PLAYWRIGHT_SSH_USER ?? 'root'

const baseArgs = [
  '-o', 'StrictHostKeyChecking=no',
  '-o', 'UserKnownHostsFile=/dev/null',
  '-o', 'LogLevel=ERROR',
  '-i', sshKey,
]

export function ssh(cmd: string, opts: { throw?: boolean } = {}): string {
  const args = [...baseArgs, `${sshUser}@${deviceHost}`, cmd]
  try {
    return execFileSync('ssh', args, { encoding: 'utf8', timeout: 120_000 })
  } catch (e: any) {
    if (opts.throw === false) {
      return (e.stdout?.toString() ?? '') + (e.stderr?.toString() ?? '')
    }
    throw e
  }
}

export function scpFrom(remote: string, local: string, opts: { throw?: boolean } = {}): void {
  const args = [...baseArgs, '-r', `${sshUser}@${deviceHost}:${remote}`, local]
  try {
    execFileSync('scp', args, { encoding: 'utf8', timeout: 120_000 })
  } catch (e) {
    if (opts.throw !== false) throw e
  }
}
