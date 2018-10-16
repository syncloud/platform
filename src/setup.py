from setuptools import setup
from os.path import join, dirname

version = open(join(dirname(__file__), 'version')).read().strip()

setup(
    name='syncloud-platform',
    version=version,
    packages=['syncloud_platform',
              'syncloud_platform.insider',
              'syncloud_platform.auth',
              'syncloud_platform.board',
              'syncloud_platform.application',
              'syncloud_platform.gaplib',
              'syncloud_platform.rest',
              'syncloud_platform.rest.facade',
              'syncloud_platform.rest.model',
              'syncloud_platform.config',
              'syncloud_platform.snap',
              'syncloud_platform.control',
              'syncloud_platform.certificate',
              'syncloud_platform.certificate.certbot',
              'syncloud_platform.disks',
              'syncloud_platform.log',
              'syncloud_platform.network'],
    namespace_packages=['syncloud_platform'],
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
