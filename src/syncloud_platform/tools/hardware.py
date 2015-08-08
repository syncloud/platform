import json
from subprocess import check_output
from syncloud_platform.systemd import systemctl


class Hardware:

    def available_disks(self, lshw_output=None, mount_output=None):
        if not lshw_output:
            lshw_output = check_output('lshw -json', shell=True)
        return self.__find_disks([], json.loads(lshw_output), mount_output)

    def __find_disks(self, acc, node, mount_output):
        if node['class'] == 'disk' and node['id'] == 'disk':
            disk = self.__parse_disk(node, mount_output)
            if disk.partitions:
                acc.append(disk)
        else:
            if 'children' in node:
                for sub_node in node['children']:
                    self.__find_disks(acc, sub_node, mount_output)
        return acc

    def __parse_disk(self, node, mount_output):
        if 'product' in node:
            name = node['product'].split(' ')[0]
        else:
            name = node['description']
        disk = Disk(name)
        for part in node['children']:
            logicalname = part['logicalname']
            if type(logicalname) is list:
                logicalname = logicalname[0]

            mountable = True
            if 'configuration' in part:
                if 'state' in part and part['state'] == 'mounted':
                    mountable = False

            if 'capabilities' in part and 'extended' in part['capabilities']:
                mountable = False

            mount_info = self.mounted_disk(logicalname, mount_output)
            mount_point = None
            if mount_info:
                mount_point = mount_info.dir
                mountable = False

            if mountable or mount_point == '/opt/disk':
                disk.partitions.append(
                    Partition(part['physid'], part['size'] / (1024 * 1024), logicalname, mount_point))
        return disk

    def mounted_disk(self, device, mount_output=None):
        if not mount_output:
            mount_output = check_output('mount', shell=True)
        for entry in mount_output.splitlines():
            if entry.startswith('{0} on'.format(device)):
                parts_on = entry.split(' on ')
                device = parts_on[0]
                parts_type = parts_on[1].split(' type ')
                dir = parts_type[0]
                parts_options = parts_type[1].split(' ')
                type = parts_options[0]
                return MountEntry(device, dir, type, parts_options[1].strip('()'))
        return None

    def activate_disk(self, device):
        systemctl.remove_mount()
        check_output('udisksctl mount -b {0}'.format(device), shell=True)
        mount_entry = self.mounted_disk(device)
        check_output('udisksctl unmount -b {0}'.format(device), shell=True)
        systemctl.add_mount(mount_entry)

    def deactivate_disk(self):
        systemctl.remove_mount()


class Partition:
    def __init__(self, id, size, device, mount_point):
        self.id = id
        self.size = size
        self.device = device
        self.mount_point = mount_point
        self.label = '{0} {1} Mb'.format(id, round(size))


class Disk:
    def __init__(self, name):
        self.partitions = []
        self.name = name


class MountEntry:

    def __init__(self, device, dir, type, options):
        self.device = device
        self.dir = dir
        self.type = type
        self.options = options
