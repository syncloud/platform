from crontab import CronTab
from syncloudlib import logger


class PlatformCron:
    def __init__(self, platform_config):
        self.platform_config = platform_config
        self.cron = CronTab(user=self.platform_config.cron_user())
        self.log = logger.get_logger('cron')

    def remove(self):
        self.log.info("remove crontab task")
        self.cron.remove_all()
        self.cron.write()

    def create(self):
        self.log.info("create crontab task")
        ci_job = self.cron.new(command=self.platform_config.cron_cmd())
        ci_job.setall(self.platform_config.cron_schedule())
        self.cron.write()
