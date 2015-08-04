import json
from subprocess import check_output


class Hardware:

    def __init__(self, test_lshw_output=None):
        self.test_lshw_output = test_lshw_output

    def disks(self):
        return self.find_disks([], json.loads(self.__read()))

    def find_disks(self, acc, node):
        if node['class'] == 'disk' and node['id'] == 'disk':
            acc.append(node)
        else:
            if 'children' in node:
                for sub_node in node['children']:
                    self.find_disks(acc, sub_node)
        return acc

    def __read(self):
        if self.test_lshw_output:
            return self.test_lshw_output
        else:
            return check_output('lshw -json', shell=True)
