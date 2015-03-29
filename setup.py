from setuptools import setup
from os.path import join, dirname

requirements = [
    'configobj==4.7.2',
    'requests==2.2.1',
    'urllib3==1.7.1',
    'IPy==0.82a',
    'convertible',
    'syncloud-app',
    'syncloud-image-tools'
]


version = open(join(dirname(__file__), 'version')).read().strip()

setup(
    name='syncloud-platform',
    version=version,
    scripts=[
        'bin/insider',
        'bin/syncloud-platform-post-install',
        'bin/syncloud-platform-pre-remove',
        'bin/syncloud-insider-post-install',
        'bin/syncloud-base-installer',
        'bin/install-java',
        'bin/syncloud-image-base-post-install',
        'bin/syncloud-link-data.sh',
        'bin/syncloud-boot-installer',
        'bin/syncloud-resize-sd',
        'bin/syncloud-resize-sd-partition',
        'bin/remote',
        'bin/syncloud-remote-pre-remove',
        'bin/syncloud-remote-post-install',
        'bin/install-avahi',
        'bin/syncloud-discovery-pre-remove'
    ],
    packages=['syncloud', 'syncloud.insider'],
    namespace_packages=['syncloud'],
    data_files=[
        ('insider/config', ['config/insider.cfg']),
        ('syncloud-image-boot/config', ['config/udisks/udisks-glue.conf']),
        ('/etc/sudoers.d', ['config/sudoers.d/www-data']),
        ('/etc/polkit-1/localauthority/50-local.d', ['config/polkit/55-storage.pkla']),
        ('/etc/udev/rules.d', ['config/udev/99-syncloud.udisks.rules']),
        ('/lib/systemd/system', ['config/systemd/ntpdate.service', 'config/systemd/udisks-glue.service']),
        ('/etc/init.d', ['bin/syncloud-resize-sd'])
    ],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
