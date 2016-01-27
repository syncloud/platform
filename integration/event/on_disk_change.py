#!/opt/app/platform/python/bin/python
from syncloud_platform.api import storage
import os


storage_dir = os.path.realpath(storage.init('platform', 'platform'))

with open('/tmp/on_disk_change.log', 'w+') as f:
    f.write(storage_dir)
