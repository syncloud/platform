from syncloud_platform.config.config import PlatformConfig, PLATFORM_CONFIG_DIR


def platform_version():
    config = PlatformConfig(PLATFORM_CONFIG_DIR)
    app_dir = config.app_dir()
    version_filename = '{0}/META/version'.format(app_dir)
    with open(version_filename) as f:
        content = f.readline()
        return content.strip()
    raise Exception("Can't get platform version")
