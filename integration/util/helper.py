import time
from integration.util.ssh import run_scp, run_ssh
from subprocess import check_output
from os.path import split
import convertible
import requests

SAM = '/opt/app/sam/bin/sam --debug'
SAM_INSTALL = '{0} install'.format(SAM)
SNAP = 'snap'
SNAP_INSTALL = '{0} install --devmode'.format(SNAP)


def local_install(host, password, app_archive_path, installer):
    _, app_archive = split(app_archive_path)
    run_scp('{0} root@{1}:/'.format(app_archive_path, host), password=password)
    cmd = SAM_INSTALL
    if installer == 'snapd':
        cmd = SNAP_INSTALL
    run_ssh(host, 'ls -la /{0}'.format(app_archive), password=password)
    run_ssh(host, '{0} /{1}'.format(cmd, app_archive), password=password)


def local_remove(host, password, installer, app):
    cmd = SAM
    if installer == 'snapd':
        cmd=SNAP
    run_ssh(host, '{0} remove {1}'.format(cmd, app), password=password)


def wait_for_platform_web(host):
    print(check_output('while ! nc -w 1 -z {0} 81; do sleep 1; done'.format(host), shell=True))
    print(check_output('while ! nc -w 1 -z {0} 80; do sleep 1; done'.format(host), shell=True))


def wait_for_sam(public_web_session, host):
    sam_running = True
    attempts = 200
    attempt = 0
    while sam_running and attempt < attempts:
        try:
            response = public_web_session.get('http://{0}/rest/settings/sam_status'.format(host))
            if response.status_code == 200:
                json = convertible.from_json(response.text)
                sam_running = json.is_running
        except Exception, e:
            print(e.message)

        print("attempt: {0}/{1}".format(attempt, attempts))
        attempt += 1
        time.sleep(1)


def wait_for_rest(public_web_session, host, url, code):
    
    attempt=0
    attempt_limit=10
    while attempt < attempt_limit:
        try:
            response = public_web_session.get('http://{0}{1}'.format(host, url))
            if response.text:
                print(response.text)
            print('code: {0}'.format(response.status_code))
            if response.status_code == code:
                return
        except Exception, e:
            print(e.message)
        time.sleep(1)
        attempt = attempt + 1



