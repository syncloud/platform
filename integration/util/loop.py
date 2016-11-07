from integration.util.ssh import run_ssh


def loop_device_cleanup(dev_file, password):
    print('cleanup')
    for loop in run_ssh('losetup -j {0} -O NAME'.format(dev_file), password=password).splitlines():
        if 'loop' in loop:
            print(loop)
            run_ssh('umount {0}'.format(loop), throw=False, password=password)
            run_ssh('losetup -d {0}'.format(loop), throw=False, password=password)

    run_ssh('losetup', password=password)

    for loop in run_ssh('dmsetup ls', password=password).splitlines():
        if 'loop0p1' in loop:
            run_ssh('sudo dmsetup remove loop0p1', password=password)
        if 'loop0p2' in loop:
            run_ssh('sudo dmsetup remove loop0p2', password=password)

    run_ssh('rm -rf {0}'.format(dev_disk), throw=False, password=password)
