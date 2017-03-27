import time
from integration.util.ssh import set_docker_ssh_port, run_scp, run_ssh
from subprocess import check_output
from os.path import split
import convertible
import requests

SAM='/opt/app/sam/bin/sam --debug'
SAM_INSTALL='{0} install'.format(SAM)
SNAP='snap'
SNAP_INSTALL='{0} install --devmode'.format(SNAP)

def local_install(password, app_archive_path, installer):
    _, app_archive = split(app_archive_path)
    run_scp('{0} root@localhost:/'.format(app_archive_path), password=password)
    cmd=SAM_INSTALL
    if installer == 'snapd':
        cmd=SNAP_INSTALL
    run_ssh('{0} /{1}'.format(cmd, app_archive), password=password)

def local_remove(password, installer, app):
    cmd=SAM
    if installer == 'snapd':
        cmd=SNAP
    run_ssh('{0} remove {1}'.format(cmd, app), password=password)


def wait_for_platform_web():
    print(check_output('while ! nc -w 1 -z localhost 81; do sleep 1; done', shell=True))
    print(check_output('while ! nc -w 1 -z localhost 80; do sleep 1; done', shell=True))


def wait_for_sam(public_web_session):
    sam_running = True
    while sam_running:
        try:
            response = public_web_session.get('http://localhost/rest/settings/sam_status')
            if response.status_code == 200:
                json = convertible.from_json(response.text)
                sam_running = json.is_running
        except Exception, e:
            print(e.message)
        time.sleep(1)


def wait_for_rest(public_web_session, url, code):
    
    attempt=0
    attempt_limit=10
    while attempt < attempt_limit:
        try:
            response = public_web_session.get('http://localhost{0}'.format(url))
            if response.text:
                print(response.text)
            print('code: {0}'.format(response.status_code))
            if response.status_code == code:
                return
        except Exception, e:
            print(e.message)
        time.sleep(1)
        attempt = attempt + 1



