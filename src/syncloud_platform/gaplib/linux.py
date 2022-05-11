from subprocess import check_output, call, CalledProcessError
import pwd
from os import environ
import massedit
from IPy import IP


def useradd(user, home_folder=None):
    try:
        pwd.getpwnam(user)
        return 'user {0} exists'.format(user)
    except KeyError:
        options = '-r -s /bin/false'
        if home_folder:
            home_folder_options = '-m -d {0}'.format(home_folder)
            options = home_folder_options + ' ' + options
        command_line = '/usr/sbin/useradd {0} {1}'.format(options, user)
        return check_output(command_line, shell=True).decode()


def fix_locale():
    if 'LANG' in environ:
        lang = environ['LANG']
        if lang not in check_output('locale -a 2>&1', shell=True).decode():
            print("generating locale: {0}".format(lang))
            __fix_locale_gen(lang)
            check_output('locale-gen')


def __fix_locale_gen(lang, locale_gen='/etc/locale.gen'):
    editor = massedit.MassEdit()
    editor.append_code_expr("re.sub('# {0}', '{0}', line)".format(lang))
    editor.edit_file(locale_gen)


def parted(device):
    return check_output('parted {0} unit % print free --script --machine'.format(device).split(" ")).decode()
