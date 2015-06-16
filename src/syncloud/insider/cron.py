import os
from os import system
from subprocess import Popen, PIPE

class Cron:

    def __init__(self, main_script_path, log_path, period_mins):
        self.log_path = log_path
        self.period_mins = period_mins
        path, script = os.path.split(main_script_path)
        self.main_script_name = script
        self.command_line = main_script_path

    def off(self):
        if not self.enabled():
            return 'already disabled'
        else:
            cmd = "crontab -l | sed '/{}/d' | crontab -".format(self.main_script_name)
            system(cmd)
            return 'disabled'

    def on(self):
        if self.enabled():
            return 'already enabled'
        else:
            cmd = "{0} sync_all > {1} 2>&1".format(self.command_line, self.log_path)
            crontab_cmd = "(crontab -l; echo \"*/{} * * * * {}\") | crontab -".format(self.period_mins, cmd)
            system(crontab_cmd)
            return 'enabled'

    def enabled(self):
        cmd = "crontab -l | grep {} | wc -l".format(self.main_script_name)
        return int(Popen(cmd, shell=True, stdout=PIPE).stdout.read()) > 0