import unittest
from subprocess import Popen, PIPE
from os import system


class TestInsider(unittest.TestCase):

    def test_find_available_port_first_gap(self):
        pre_test_port_count = int(Popen("../bin/insider list_ports | wc -l", shell=True, stdout=PIPE).stdout.read())
        system("../bin/insider add_port 11000")
        system("../bin/insider add_port 11001")
        system("../bin/insider list_ports")
        post_test_port_count = int(Popen("../bin/insider list_ports | wc -l", shell=True, stdout=PIPE).stdout.read())
        system("../bin/insider remove_port 11000")
        system("../bin/insider remove_port 11001")
        system("../bin/insider list_ports")

        self.assertEquals(pre_test_port_count + 2, post_test_port_count)