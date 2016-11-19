import time
from integration.util.ssh import set_docker_ssh_port, run_scp, run_ssh
from subprocess import check_output
from os.path import split
import convertible
import requests


def local_install(password, app_archive_path):
    _, app_archive = split(app_archive_path)
    run_scp('{0} root@localhost:/'.format(app_archive_path), password=password)
    run_ssh('/opt/app/sam/bin/sam --debug install /{0}'.format(app_archive), password=password)
    set_docker_ssh_port(password)
    run_ssh("sed -i 's/certbot_test_cert.*/certbot_test_cert: true/g' /opt/data/platform/config/platform.cfg ",
            password=password)
    run_ssh('systemctl restart platform-uwsgi-public', password=password)

    time.sleep(3)


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
