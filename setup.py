from setuptools import setup
from os.path import join, dirname
from sys import exec_prefix

# Use prefix for virtual env
prefix = ''
if not exec_prefix == '/usr':
    prefix = join(exec_prefix, 'local')

requirements = [
    'configobj==4.7.2',
    'requests==2.2.1',
    'urllib3==1.7.1',
    'IPy==0.82a',
    'beautifulsoup4==4.3.2',
    'convertible==0.13',
    'Flask==0.10.1',
    'psutil==2.1.3'
]


version = open(join(dirname(__file__), 'version')).read().strip()

setup(
    name='syncloud-platform',
    version=version,
    scripts=[
        'bin/insider',
        'bin/syncloud-platform-post-install',
        'bin/syncloud-platform-post-upgrade',
        'bin/syncloud-platform-pre-remove',
        'bin/syncloud-base-installer',
        'bin/install-java',
        'bin/syncloud-link-data.sh',
        'bin/syncloud-boot-installer',
        'bin/syncloud-resize-sd',
        'bin/syncloud-resize-sd-partition',
        'bin/remote',
        'bin/syncloud-remote-pre-remove',
        'bin/syncloud-remote-post-install',
        'bin/install-avahi',
        'bin/syncloud-discovery-pre-remove',
        'bin/syncloud-apache-post-install',
        'bin/syncloud-apache',
        'bin/syncloud-cli',
        'bin/syncloud-server-post-upgrade',
        'bin/sam',
        'bin/syncloud-id',
        'bin/syncloud-ping',
        'bin/syncloud_ssh_keys_check',
        'bin/syncloud_ssh_keys_generate',
        'bin/cpu_frequency'
    ],
    packages=['syncloud', 'syncloud.insider', 'syncloud.remote', 'syncloud.apache', 'syncloud.server',
              'syncloud.sam', 'syncloud.app', 'syncloud.tools', 'syncloud.tools.cpu'],
    namespace_packages=['syncloud'],
    data_files=[
        ('insider/config', ['config/insider.cfg']),
        ('syncloud-image-boot/config', ['config/udisks/udisks-glue.conf']),
        (prefix + '/etc/sudoers.d', ['config/sudoers.d/www-data']),
        (prefix + '/etc/polkit-1/localauthority/50-local.d', ['config/polkit/55-storage.pkla']),
        (prefix + '/etc/udev/rules.d', ['config/udev/99-syncloud.udisks.rules']),
        (prefix + '/lib/systemd/system', [
            'config/systemd/ntpdate.service',
            'config/systemd/udisks-glue.service',
            'config/systemd/syncloud-resize-sd.service',
            'config/systemd/cpu-frequency.service',
            'config/systemd/insider-sync.service']),
        ('syncloud-apache/config', [
            'config/http.conf',
            'config/https.conf']),
        ('syncloud-server/config', ['config/server.wsgi']),
        ('syncloud-server/apache', [
            'apache/syncloud-server-http.conf',
            'apache/syncloud-server-https.conf']),
        (prefix + '/var/www/syncloud-server', [
            'www/favicon.ico',
            'www/index.html']),
        (prefix + '/var/www/syncloud-server/images', [
            'www/images/image-ci-128.png',
            'www/images/owncloud-128.png']),
        ('sam/config', ['config/sam.cfg'])
    ],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
