import { execFileSync } from 'node:child_process'

export const deviceHost = process.env.PLAYWRIGHT_DEVICE_HOST ?? 'bookworm-amd64'
export const sshUser = process.env.PLAYWRIGHT_SSH_USER ?? 'root'
export const sshPassword = process.env.PLAYWRIGHT_SSH_PASSWORD ?? 'Password1'

const baseArgs = [
  '-o', 'StrictHostKeyChecking=no',
  '-o', 'UserKnownHostsFile=/dev/null',
  '-o', 'LogLevel=ERROR',
]

export function ssh(cmd: string, opts: { throw?: boolean } = {}): string {
  const args = ['-p', sshPassword, 'ssh', ...baseArgs, `${sshUser}@${deviceHost}`, cmd]
  try {
    return execFileSync('sshpass', args, { encoding: 'utf8', timeout: 120_000 })
  } catch (e: any) {
    if (opts.throw === false) {
      return (e.stdout?.toString() ?? '') + (e.stderr?.toString() ?? '')
    }
    throw e
  }
}

export function scpFrom(remote: string, local: string, opts: { throw?: boolean } = {}): void {
  const args = ['-p', sshPassword, 'scp', ...baseArgs, '-r', `${sshUser}@${deviceHost}:${remote}`, local]
  try {
    execFileSync('sshpass', args, { encoding: 'utf8', timeout: 120_000 })
  } catch (e) {
    if (opts.throw !== false) throw e
  }
}
