from os.path import join

from syncloud_platform.gaplib import gen
from syncloudlib import logger


class Nginx:
    def __init__(self, platform_config, device_info):
        self.config = platform_config
        self.device_info = device_info
        self.log = logger.get_logger('nginx')

    def init_config(self):
        domain = self.device_info.domain()
        variables = {'domain': domain.replace(".", "\\.")}
        gen.generate_file_jinja(
            join(self.config.config_dir(), 'nginx', 'public.conf'), 
            join(self.config.data_dir(), 'nginx.conf'),
            variables)
