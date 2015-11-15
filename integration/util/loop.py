from integration.util.ssh import run_ssh


def loop_device_cleanup():
    print('cleanup')
    for mount in run_ssh('mount', debug=False).splitlines():
        if 'loop' in mount:
            print(mount)

    for loop in run_ssh('losetup').splitlines():
        if 'loop0p1' in loop:
            run_ssh('losetup -d /dev/loop0', throw=False)

    run_ssh('losetup')

    for loop in run_ssh('dmsetup ls').splitlines():
        if 'loop0p1' in loop:
            run_ssh('sudo dmsetup remove loop0p1')

    for loop_disk in run_ssh('ls -la /tmp').splitlines():
        if '/tmp/disk' in loop_disk:
            run_ssh('rm -fr /tmp/disk', throw=False)
