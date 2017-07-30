from subprocess import check_output, STDOUT, CalledProcessError

import time


def run_scp(command, throw=True, debug=True, password='syncloud'):
    return _run_command('scp -o StrictHostKeyChecking=no {0}'.format(command), throw, debug, password)


def run_ssh(host, command, throw=True, debug=True, password='syncloud', retries=0, sleep=1, env_vars=''):
    retry = 0
    while True:
        try:
            command='{0}{1}'.format(env_vars, command)
            return _run_command('ssh -o StrictHostKeyChecking=no root@{0} "{1}"'.format(host, command), throw, debug, password)
        except Exception, e:
            if retry >= retries:
                raise e
            retry += 1
            time.sleep(sleep)
            print('retrying {0}'.format(retry))


def ssh_command(password, command):
    return 'sshpass -p {0} {1}'.format(password, command)


def _run_command(command, throw, debug, password):
    try:
        print('ssh command: {0}'.format(command))
        output = check_output(ssh_command(password, command), shell=True, stderr=STDOUT).strip()
        if debug:
            print output
            print
        return output
    except CalledProcessError, e:
        print(e.output)
        if throw:
            raise e
