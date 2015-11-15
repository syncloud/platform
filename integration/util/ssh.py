from subprocess import check_output


DOCKER_SSH_PORT = 2222
SSH = 'sshpass -p syncloud ssh -o StrictHostKeyChecking=no -p {0} root@localhost'.format(DOCKER_SSH_PORT)


def set_docker_ssh_port():
    run_ssh("sed -i 's/ssh_port.*/ssh_port:{0}/g' /opt/app/platform/config/platform.cfg".format(DOCKER_SSH_PORT))


def run_ssh(command, throw=False, debug=True):
    try:
        output = check_output('{0} {1}'.format(SSH, command), shell=True).strip()
        if debug:
            print('ssh command: {0}'.format(command))
            print output
            print
        return output
    except Exception, e:
        print(e.message)
        if throw:
            raise e
