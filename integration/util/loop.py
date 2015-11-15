from integration.util.ssh import run_ssh


def loop_device_cleanup():
    print('cleanup')
    for mount in run_ssh('mount', debug=False).splitlines():
        if 'loop' in mount:
            print(mount)
    run_ssh('losetup -a')
    run_ssh('losetup -d /dev/loop0', throw=False)
    run_ssh('losetup -a')
    run_ssh('rm -fr /tmp/disk', throw=False)
