class Versions:

    def __init__(self, platform_config):
        self.platform_config = platform_config

    def platform_version(self):
        app_dir = self.platform_config.app_dir()
        version_filename = '{0}/META/version'.format(app_dir)
        with open(version_filename) as f:
            content = f.readline()
            return content.strip()
