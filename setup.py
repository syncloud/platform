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
    'psutil==2.1.3',
    'python-crontab==1.7.2',
    'wget==2.2',
    'massedit==0.66',
    'python-ldap==2.4.19',
    'flask-login==0.2.10'
]


version = open(join(dirname(__file__), 'version')).read().strip()

setup(
    name='syncloud-platform',
    version=version,
    scripts=[
        'bin/insider',
        'bin/syncloud-platform-post-install',
        'bin/syncloud-platform-pre-remove',
        'bin/syncloud-boot-installer',
        'bin/syncloud-link-data.sh',
        'bin/syncloud-resize-sd',
        'bin/syncloud-resize-sd-partition',
        'bin/syncloud-cli',
        'bin/sam',
        'bin/syncloud-id',
        'bin/syncloud-ping',
        'bin/cpu_frequency'
    ],
    packages=['syncloud', 'syncloud.insider', 'syncloud.server',
              'syncloud.sam', 'syncloud.app', 'syncloud.tools', 'syncloud.tools.cpu', 'syncloud.systemd',
              'syncloud.server.rest', 'syncloud.config', 'syncloud.installer'],
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
        ('sam/config', ['config/sam.cfg'])
    ],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
