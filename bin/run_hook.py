from syncloud.systemd.systemctl import add_service, remove_service
from syncloud.tools import app
from syncloud.tools.nginx import Nginx

import sys

script_filename = sys.argv[1]
execfile(script_filename)