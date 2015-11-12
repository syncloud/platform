#!/opt/app/platform/python/bin/python
from syncloud_platform.api import storage

storage_dir = storage.init('platform', 'platform')

with open('/tmp/on_disk_change.log', 'w+') as f:
    f.write(storage_dir)
