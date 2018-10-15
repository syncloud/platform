from syncloud_platform.snap.models import App, AppVersions
from syncloud_platform.snap.snap import join_apps


def test_join_apps():

    installed_app1 = App()
    installed_app1.id = 'id1'
    installed_app_version1 = AppVersions()
    installed_app_version1.installed_version = 'v1'
    installed_app_version1.current_version = None
    installed_app_version1.app = installed_app1

    installed_app2 = App()
    installed_app2.id = 'id2'
    installed_app_version2 = AppVersions()
    installed_app_version2.installed_version = 'v1'
    installed_app_version2.current_version = None
    installed_app_version2.app = installed_app2

    installed_apps = [installed_app_version1, installed_app_version2]

    store_app2 = App()
    store_app2.id = 'id2'
    store_app_version2 = AppVersions()
    store_app_version2.installed_version = None
    store_app_version2.current_version = 'v2'
    store_app_version2.app = store_app2

    store_app3 = App()
    store_app3.id = 'id3'
    store_app_version3 = AppVersions()
    store_app_version3.installed_version = None
    store_app_version3.current_version = 'v2'
    store_app_version3.app = store_app3

    store_apps = [store_app_version2, store_app_version3]

    all_apps = sorted(join_apps(installed_apps, store_apps), key=lambda app: app.app.id)

    assert len(all_apps) == 3

    assert all_apps[0].app.id == 'id1'
    assert all_apps[0].installed_version == 'v1'
    assert all_apps[0].current_version is None

    assert all_apps[1].app.id == 'id2'
    assert all_apps[1].installed_version == 'v1'
    assert all_apps[1].current_version == 'v2'

    assert all_apps[2].app.id == 'id3'
    assert all_apps[2].installed_version is None
    assert all_apps[2].current_version == 'v2'
