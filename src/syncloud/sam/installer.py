import os
from os.path import isfile
from subprocess import check_output
import pwd
from syncloud.app import logger
from syncloud.sam.manager import get_sam
import wget
import tarfile
import massedit


class Installer:
    def __init__(self):
        self.log = logger.get_logger('sam.installer')
        self.sam = get_sam()

    def install(self, app_id, from_file=None, owner=None, owner_home=None, apps_root='/opt/app'):

        lang = os.environ['LANG']
        if lang not in check_output(['locale', '-a']):
            print("generating locale: {0}".format(lang))
            fix_locale_gen(lang)
            check_output('locale-gen')

        if not from_file:
            archive = '{0}.tar.gz'.format(app_id)
            arch = check_output('uname -m', shell=True).strip()
            url = 'http://apps.syncloud.org/{0}/{1}/{2}'.format(self.sam.get_release(), arch, archive)

            self.log.info('installing from: {0}'.format(url))

            temp = '/tmp/{0}'.format(archive)
            if isfile(temp):
                os.remove(temp)

            self.log.info("saving {0} to {1}".format(url, temp))
            filename = wget.download(url, temp)
        else:
            filename = from_file
            self.log.info('installing from: {0}'.format(from_file))

        self.log.info("extracting {0}".format(filename))
        tarfile.open(filename).extractall(apps_root)
        if owner:
            try:
                pwd.getpwnam(owner)
            except KeyError:
                if owner_home:
                    owner_home = '--home {0}'.format(owner_home)
                else:
                    owner_home = ''
                self.log.info(check_output('/usr/sbin/useradd -r -s /bin/false {0} {1}'.format(owner, owner_home), shell=True))
            self.log.info(check_output('chown -R {0}. {1}/{2}'.format(owner, apps_root, app_id), shell=True))


def fix_locale_gen(lang, locale_gen='/etc/locale.gen'):
    massedit.edit_files([locale_gen], ["re.sub('# {0}', '{0}', line)".format(lang)], dry_run=False)
