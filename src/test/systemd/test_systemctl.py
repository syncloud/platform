from syncloud_platform.systemd.systemctl import __dir_to_systemd_mount_filename


def test_dir_to_systemd_mount_filename():
    assert __dir_to_systemd_mount_filename('/dir1/dir2') == 'dir1-dir2.mount'