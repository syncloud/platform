from crontab import CronTab
from syncloud_platform.config.config import PlatformConfig

class PlatformCron:

    def __init__(self, platform_config=None):
        if not platform_config:
            platform_config = PlatformConfig()
        self.platform_config = platform_config
        self.cron = CronTab(user=self.platform_config.cron_user())

    def remove(self):
        print("remove crontab task")
        for job in self.cron.find_command(self.platform_config.cron_user()):
            self.cron.remove(job)
        self.cron.write()

    def create(self):
        print("create crontab task")
        ci_job = self.cron.new(command=self.platform_config.cron_cmd())
        ci_job.setall(self.platform_config.cron_schedule())
        self.cron.write()
