import os
import subprocess
import signal
import shutil
import tempfile


class FakePyPi:

    def __init__(self):
        self.port = 10003
        self.repo_path = tempfile.mkdtemp()

    def add_app_version(self, name, file_path):
        app_path = os.path.join(self.repo_path, name)
        if not os.path.exists(app_path):
            os.makedirs(app_path)
        path, name = os.path.split(file_path)
        app_version_path = os.path.join(app_path, name)
        shutil.copyfile(file_path, app_version_path)

    def start(self):
        cmd = 'pypi-server -p {} {}'.format(self.port, self.repo_path)
        self.server = subprocess.Popen(cmd, shell=True, preexec_fn=os.setsid)

    def stop(self):
        os.killpg(self.server.pid, signal.SIGKILL)
        os.waitpid(self.server.pid, 0)
        self.clear()

    def clear(self):
        if os.path.isdir(self.repo_path):
            shutil.rmtree(self.repo_path)

    def url(self):
        return 'http://localhost:{}/simple/'.format(self.port)
