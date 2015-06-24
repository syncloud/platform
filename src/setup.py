from setuptools import setup
from os.path import join, dirname

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
        'bin/syncloud-platform-post-install',
        'bin/syncloud-platform-pre-remove',
        'bin/sam'
    ],
    packages=['syncloud', 'syncloud.insider', 'syncloud.server',
              'syncloud.app', 'syncloud.tools', 'syncloud.tools.cpu', 'syncloud.systemd',
              'syncloud.server.rest', 'syncloud.config', 'syncloud.installer'],
    namespace_packages=['syncloud'],
    data_files=[
        ('sam/config', ['config/sam.cfg'])
    ],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
