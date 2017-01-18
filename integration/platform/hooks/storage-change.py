#!/opt/app/platform/python/bin/python
from syncloud_platform.application import api
import os

app = api.get_app_setup('platform')
storage_dir = os.path.realpath(app.init_storage('platform'))

with open('/tmp/on_disk_change.log', 'w+') as f:
    f.write(storage_dir)
