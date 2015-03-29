from os.path import join
from syncloud.tools.facade import Facade

tools_facade = Facade()
config_dir = join(tools_facade.usr_local_dir(), 'syncloud-apache', 'config')

apache_root = '/etc/apache2'
sites_available_dir = join(apache_root, 'sites-available')
sites_enabled_dir = join(apache_root, 'sites-enabled')

http_file_template = join(config_dir, 'http.conf')
http_file = "{}/http".format(sites_available_dir)

https_file_template = join(config_dir, 'https.conf')
https_file = "{}/https".format(sites_available_dir)

http_web_root = '/var/www/http'
https_web_root = '/var/www/https'

log_dir = '/var/log/apache2'

cert_file = "/etc/ssl/certs/syncloud.crt"
key_file = "/etc/ssl/private/syncloud.key"

http_include_dir = "conf.http.d"
http_include_dir_full = join(apache_root, http_include_dir)

https_include_dir = "conf.https.d"
https_include_dir_full = join(apache_root, https_include_dir)
