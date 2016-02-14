from setuptools import setup
from os.path import join, dirname

requirements = [
    'urllib3==1.10.4',
    'requests==2.7.0',
    'beautifulsoup4==4.3.2',
    'IPy==0.83',
    'python-crontab==1.7.2',
    'massedit==0.66',
    'Flask==0.10.1',
    'flask-login==0.2.10',
    'syncloud-lib==2'
]


version = open(join(dirname(__file__), 'version')).read().strip()

setup(
    name='syncloud-platform',
    version=version,
    packages=['syncloud_platform', 'syncloud_platform.insider', 'syncloud_platform.api', 'syncloud_platform.auth',
              'syncloud_platform.tools', 'syncloud_platform.tools.cpu', 'syncloud_platform.systemd',
              'syncloud_platform.rest', 'syncloud_platform.config', 'syncloud_platform.sam',
              'syncloud_platform.rest.facade', 'syncloud_platform.rest.model', 'syncloud_platform.di',
              'syncloud_platform.tools.disk', 'syncloud_platform.log'],
    namespace_packages=['syncloud_platform'],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
