from subprocess import check_output
import pwd
from os import environ
import massedit


def useradd(user):
    try:
        pwd.getpwnam(user)
        return 'user {0} exists'.format(user)
    except KeyError:
        return check_output('/usr/sbin/useradd -r -s /bin/false {0}'.format(user), shell=True)


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
    ip = check_output(["hostname", "-I"]).split(" ")[0]
    if not ip:
        raise(Exception("Can't get local ip address"))
    return ip
