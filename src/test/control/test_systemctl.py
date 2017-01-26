from syncloud_platform.control.systemctl import dir_to_systemd_mount_filename


def test_dir_to_systemd_mount_filename():
    assert dir_to_systemd_mount_filename('/dir1/dir2') == 'dir1-dir2.mount'
