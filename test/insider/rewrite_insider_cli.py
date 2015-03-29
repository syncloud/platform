from ConfigParser import ConfigParser
import unittest
from subprocess import check_output, call, CalledProcessError
from os.path import dirname, join

test_dir = dirname(__file__)
root_dir = join(dirname(__file__), '..')
cmd = '{}/bin/insider --config-path={}/test-config/ --debug'.format(root_dir, root_dir)


class TestCli(unittest.TestCase):

    def setUp(self):
        clean_up()
        check_output('mkdir {}/test-config'.format(root_dir), shell=True)
        copy_conf('insider.cfg')

    def test_user_domain_not_activated(self):

        self.assertEquals(1, call('{} user_domain'.format(cmd), shell=True))

    def test_full_name_activated(self):

        copy_conf('domain.json')
        output = check_output('{} user_domain'.format(cmd), shell=True)
        self.assertTrue("\"data\": \"domain\"" in output, output)

    def test_user_domain_activated(self):

        copy_conf('domain.json')
        output = check_output('{} user_domain'.format(cmd), shell=True)
        self.assertTrue("\"data\": \"domain\"" in output, output)

    def test_cron(self):

        output = check_output('{} cron_on'.format(cmd), shell=True)
        self.assertTrue("\"enabled\"" in output, output)
        output = check_output('{} cron_on'.format(cmd), shell=True)
        self.assertTrue("\"already enabled\"" in output)

        output = check_output('{} cron_off'.format(cmd), shell=True)
        self.assertTrue("\"disabled\"" in output)
        output = check_output('{} cron_off'.format(cmd), shell=True)
        self.assertTrue("\"already disabled\"" in output)

    def test_set_redirect_info(self):

        copy_conf('domain.json')
        call('{} set_redirect_info domain1 http://api.domain1'.format(cmd), shell=True)
        parser = ConfigParser()
        parser.read('{}/test-config/insider.cfg'.format(root_dir))

        self.assertEquals(parser.get('redirect', 'domain'), 'domain1')
        self.assertEquals(parser.get('redirect', 'api_url'), 'http://api.domain1')

        output = check_output('{} user_domain'.format(cmd), shell=True)
        self.assertTrue("\"domain\"" in output)

        call('{} set_redirect_info domain2 http://api.domain2'.format(cmd), shell=True)
        parser = ConfigParser()
        parser.read('{}/test-config/insider.cfg'.format(root_dir))

        self.assertEquals(parser.get('redirect', 'domain'), 'domain2')
        self.assertEquals(parser.get('redirect', 'api_url'), 'http://api.domain2')

        output = check_output('{} user_domain'.format(cmd), shell=True)
        self.assertTrue("\"domain\"" in output)

    def test_acquire_domain(self):

        # TODO: Implement dummy web-server
        call('{} acquire_domain email pass domain'.format(cmd), shell=True)

    def test_drop_domain(self):

        call('{} drop_domain'.format(cmd), shell=True)
        self.assertEquals(1, call('{} user_domain'.format(cmd), shell=True))

    def test_service_info(self):

        copy_conf('domain.json')
        copy_conf('ports.json')
        copy_conf('services.json')

        output = check_output('{} service_info ssh'.format(cmd), shell=True)
        self.assertTrue("\"type\": \"_ssh._tcp\"" in output, output)
        self.assertTrue("\"port\": \"1022\"" in output, output)
        self.assertTrue("\"external_port\": 10000" in output, output)
        self.assertTrue("\"external_host\": \"device.domain.localhost\"" in output, output)

    def test_service_info_not_found(self):

        copy_conf('domain.json')
        copy_conf('ports.json')
        copy_conf('services.json')

        try:
            output = check_output('{} service_info owncloud'.format(cmd), shell=True)
            self.fail("should no succeed")
        except CalledProcessError, e:
            self.assertTrue("not found" in e.output, e.output)

    def tearDown(self):
        clean_up()


def copy_conf(conf):
    call('cp {}/conf/{} {}/test-config'.format(test_dir, conf, root_dir), shell=True)


def clean_up():
    check_output('rm -rf {}/test-config'.format(root_dir), shell=True)