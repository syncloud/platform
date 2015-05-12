import logging

import os
from os.path import dirname, join, split

import shutil
import tempfile
import pytest
from syncloud.app import logger

from syncloud.sam.manager import get_sam
from syncloud.sam.config import SamConfig
from syncloud.sam.pip import Pip

from subprocess import check_output
from fakepypi import FakePyPi
import responses

test_dir = dirname(__file__)
logger.init(logging.DEBUG, console=True)


def text_file(path, filename, text=''):
    app_path = join(path, filename)
    f = open(app_path, 'w')
    f.write(text)
    f.close()
    return app_path


def create_app_version(name, version, pre_remove=None, post_install=None, reconfigure=None):
    temp_folder = tempfile.mkdtemp()

    app_script_content='''#!/usr/bin/env python
print("{}")'''.format(version)
    text_file(temp_folder, name, app_script_content)
    scripts = [name]

    if pre_remove:
        hook_filename = '{}-pre-remove'.format(name)
        text_file(temp_folder, hook_filename, pre_remove)
        scripts.append(hook_filename)

    if post_install:
        hook_filename = '{}-post-install'.format(name)
        text_file(temp_folder, hook_filename, post_install)
        scripts.append(hook_filename)

    if reconfigure:
        hook_filename = '{}-reconfigure'.format(name)
        text_file(temp_folder, hook_filename, reconfigure)
        scripts.append(hook_filename)

    scripts_line = ', '.join(["'"+s+"'" for s in scripts])

    setup_content='''from setuptools import setup

setup(
    name='{name}',
    version='{version}',
    scripts=[{scripts}])'''.format(name=name, version=version, scripts=scripts_line)

    text_file(temp_folder, 'setup.py', setup_content)

    check_output(['python', 'setup.py', 'sdist'], cwd=temp_folder)

    dist_filename = '{}-{}.tar.gz'.format(name, version)
    return join(temp_folder, 'dist', dist_filename)


def create_release(release, index, versions=None):
    prepare_dir_root = tempfile.mkdtemp()
    prepare_dir = join(prepare_dir_root, 'app-' + release)
    os.makedirs(prepare_dir)

    text_file(prepare_dir, 'index', index)
    if versions:
        text_file(prepare_dir, 'versions', versions)

    archive = shutil.make_archive(release, 'zip', prepare_dir_root)

    return archive


def assert_single_application(applications, id, name, current_version, installed_version):
    assert applications is not None
    assert len(applications) == 1

    test_app = applications[0]
    assert test_app.app.id == id
    assert test_app.app.name == name
    assert test_app.current_version == current_version
    assert test_app.installed_version == installed_version


def one_app_index(required=False):
    app_index_template = '''{
      "apps" : [
        {
          "name" : "test app",
          "id" : "test-app",
          "type": "admin",
          "required": %s
        }
      ]
    }'''
    return app_index_template % str(required).lower()


class BaseTest:
    def setup(self):
        self.pypi = FakePyPi()
        self.pypi.start()

        self.config_dir = tempfile.mkdtemp()
        apps_dir = tempfile.mkdtemp()
        status_dir = tempfile.mkdtemp()
        self.releases_dir = tempfile.mkdtemp()

        apps_url_template = 'file://'+self.releases_dir+'/{}.zip'

        self.config = SamConfig(join(self.config_dir, 'sam.cfg'))
        self.config.set_apps_dir(apps_dir)
        self.config.set_status_dir(status_dir)
        self.config.set_apps_url_template(apps_url_template)
        self.config.set_bin_dir('/usr/local/bin')

        self.sam = get_sam(self.config_dir, self.pypi.url())

        self.pip = Pip(self.pypi.url(), raise_on_error=False)
        self.pip.uninstall('test-app')

    def teardown(self):
        self.pypi.stop()
        self.pip.uninstall('test-app')

    def create_release(self, release, index, versions=None):
        archive_path = create_release(release, index, versions)
        _, name = split(archive_path)
        release_path = join(self.releases_dir, name)
        shutil.move(archive_path, release_path)

    def create_app_version(self, name, version, pre_remove=None, post_install=None, reconfigure=None):
        dist_path = create_app_version(name, version, pre_remove, post_install, reconfigure)
        self.pypi.add_app_version(name, dist_path)


class TestBasic(BaseTest):

    def test_list(self):
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')
        self.sam.update('release-1.0')
        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)

    def test_install(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')
        self.sam.install('test-app')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', '1.0')

        output = check_output('test-app', shell=True)
        assert output.strip() == '1.0'

    def test_upgrade(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')
        self.sam.install('test-app')

        self.create_app_version('test-app', '1.1')
        self.create_release('release-1.1', one_app_index(), 'test-app=1.1')
        self.sam.update('release-1.1')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.1', '1.0')

        self.sam.upgrade('test-app')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.1', '1.1')

        output = check_output('test-app', shell=True)
        assert output.strip() == '1.1'

    def test_remove(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')
        self.sam.install('test-app')

        self.sam.remove('test-app')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)


class TestUpdates(BaseTest):

    def test_update_simple(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)

    def test_update_same_release(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        self.create_app_version('test-app', '1.0.1')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0.1')

        self.sam.update()

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0.1', None)

    def test_update_new_release(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        self.create_app_version('test-app', '1.1')
        self.create_release('release-1.1', one_app_index(), 'test-app=1.1')

        self.sam.update('release-1.1')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.1', None)

    def test_update_bootstrap_no_apps_dir(self):
        self.config.set_apps_dir(join(self.config.apps_dir(), 'non_existent'))

        sam = get_sam(self.config_dir)

        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')
        sam.update('release-1.0')

        applications = sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)

    @responses.activate
    def test_update_no_versions(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), versions=None)

        responses.add(responses.GET,
                      "https://pypi.python.org/pypi/test-app/json",
                      status=200,
                      body='{"info": {"version": "1.0"}}',
                      content_type="application/json")

        self.sam.update('release-1.0')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)


class TestUpgradeAll(BaseTest):

    def test_update_usual_app(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')
        self.sam.update('release-1.0')
        self.sam.install('test-app')

        self.create_app_version('test-app', '1.1')
        self.create_release('release-1.1', one_app_index(), 'test-app=1.1')
        self.sam.update('release-1.1')
        self.sam.upgrade_all()

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.1', '1.1')

    def test_update_required_app(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(required=True), 'test-app=1.0')

        self.sam.update('release-1.0')
        self.sam.upgrade_all()

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', '1.0')


class Test_Reconfigure_Installed_Apps(BaseTest):

    def test_reconfigure_bad(self):
        reconfigure_content = '#!/bin/sh\nexit 1'
        self.create_app_version('test-app', '1.1', reconfigure=reconfigure_content)
        self.create_release('release-1.1', one_app_index(), 'test-app=1.1')
        self.sam.update('release-1.1')
        self.sam.install('test-app')

        with pytest.raises(Exception):
            self.sam.reconfigure_installed_apps()

    def test_reconfigure_good(self):
        reconfigure_content = '#!/bin/sh\nexit 0'
        self.create_app_version('test-app', '1.0', reconfigure=reconfigure_content)
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')
        self.sam.update('release-1.0')
        self.sam.install('test-app')

        self.sam.reconfigure_installed_apps()


class TestHooks(BaseTest):

    def test_remove_hook_good(self):
        pre_remove_content = '#!/bin/sh\nexit 0'

        self.create_app_version('test-app', '1.0', pre_remove=pre_remove_content)
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        self.sam.install('test-app')

        self.sam.remove('test-app')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)

    def test_remove_hook_missing(self):
        self.create_app_version('test-app', '1.0')
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')
        self.sam.update('release-1.0')

        self.sam.install('test-app')

        self.sam.remove('test-app')

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)

    def _test_remove_hook_bad(self):
        pre_remove_content = '#!/bin/sh\nexit 1'

        self.create_app_version('test-app', '1.0', pre_remove=pre_remove_content)
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        self.sam.install('test-app')

        try:
            self.sam.remove('test-app')
            assert False
        except Exception, e:
            assert True

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', '1.0')

    def _test_install_hook_bad(self):
        post_install_content = '#!/bin/sh\nexit 1'

        self.create_app_version('test-app', '1.0', post_install=post_install_content)
        self.create_release('release-1.0', one_app_index(), 'test-app=1.0')

        self.sam.update('release-1.0')

        try:
            self.sam.install('test-app')
            assert False
        except Exception, e:
            assert True

        applications = self.sam.list()
        assert_single_application(applications, 'test-app', 'test app', '1.0', None)
