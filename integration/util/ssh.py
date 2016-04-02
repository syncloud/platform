from subprocess import check_output, STDOUT, CalledProcessError

import time

DOCKER_SSH_PORT = 2222
SSH = 'ssh -o StrictHostKeyChecking=no -p {0} root@localhost'.format(DOCKER_SSH_PORT)
SCP = 'scp -o StrictHostKeyChecking=no -P {0}'.format(DOCKER_SSH_PORT)


def set_docker_ssh_port(password):
    run_ssh("sed -i 's/ssh_port.*/ssh_port = {0}/g' /opt/app/platform/config/platform.cfg".format(DOCKER_SSH_PORT),
            password=password)


def run_scp(command, throw=True, debug=True, password='syncloud'):
    return _run_command('{0} {1}'.format(SCP, command), throw, debug, password)


def run_ssh(command, throw=True, debug=True, password='syncloud', retries=0, sleep=1):
    retry = 0
    while True:
        try:
            return _run_command('{0} "{1}"'.format(SSH, command), throw, debug, password)
        except Exception, e:
            if retry >= retries:
                raise e
            retry = retry + 1
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
