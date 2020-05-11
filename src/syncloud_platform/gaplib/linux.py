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
        return check_output(command_line, shell=True)


def fix_locale():
    if 'LANG' in environ:
        lang = environ['LANG']
        if lang not in check_output('locale -a 2>&1', shell=True):
            print("generating locale: {0}".format(lang))
            __fix_locale_gen(lang)
            check_output('locale-gen')


def __fix_locale_gen(lang, locale_gen='/etc/locale.gen'):
    editor = massedit.MassEdit()
    editor.append_code_expr("re.sub('# {0}', '{0}', line)".format(lang))
    editor.edit_file(locale_gen)


def local_ip():
    ip = check_output(["hostname", "-I"]).decode().split(" ")[0]
    if not ip:
        raise(Exception("Can't get local ip address"))
    return ip

def local_ip_v6():
    try:
        return check_output("/snap/platform/current/bin/cli ipv6", shell=True)
    except CalledProcessError as e:
        return None

def public_ip_v4():
    try:
        return check_output("/snap/platform/current/bin/cli ipv4 piblic", shell=True)
    except CalledProcessError as e:
        return None

def is_ip_public(ip):
    return ip_type(ip) == 'PUBLIC'

def ip_type(ip):
    return IP(ip).iptype()

def parted(device):
    return check_output('parted {0} unit % print free --script --machine'.format(device).split(" "))
