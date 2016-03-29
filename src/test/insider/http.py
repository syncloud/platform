import subprocess
import os
import signal
import requests
import time


def wait_http(ip, port, status_code, timeout):
    start = time.time()
    while True:
        try:
            r = requests.get('http://{ip}:{port}'.format(ip=ip, port=port), timeout=timeout)
            if r.status_code == status_code:
                return r
        except Exception as ex:
            pass
        elapsed = time.time() - start
        if elapsed > timeout:
            raise Exception('Timeout of {} seconds happened {} seconds elapsed'.format(timeout, elapsed))


def wait_http_cant_connect(ip, port, timeout):
    start = time.time()
    while True:
        try:
            requests.get('http://{ip}:{port}'.format(ip=ip, port=port), timeout=timeout)
        except Exception as ex:
            print(ex)
            return ex
        print('waiting')
        elapsed = time.time() - start
        if elapsed > timeout:
            raise Exception('Timeout of {} seconds happened {} seconds elapsed'.format(timeout, elapsed))


class SomeHttpServer:
    def __init__(self, port):
        self.port = port

    def start(self):
        cmd = 'python -m SimpleHTTPServer {port}'.format(port=self.port)
        self.server = subprocess.Popen(cmd, shell=True, preexec_fn=os.setsid)
        wait_http('localhost', self.port, 200, timeout=5)

    def stop(self):
        os.killpg(self.server.pid, signal.SIGKILL)
        os.waitpid(self.server.pid, 0)