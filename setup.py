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
        'bin/syncloud-insider-post-install',
        'bin/syncloud-base-installer',
        'bin/install-java',
        'bin/syncloud-image-base-post-install'
    ],
    packages=['syncloud', 'syncloud.insider'],
    namespace_packages=['syncloud'],
    data_files=[('insider/config', ['config/insider.cfg'])],
    install_requires=requirements,
    description='Syncloud platform',
    long_description='Syncloud platform',
    license='GPLv3',
    author='Syncloud',
    author_email='syncloud@googlegroups.com',
    url='https://github.com/syncloud/platform')
