#!/opt/app/platform/python/bin/python
from syncloud_platform.api import info

domain = info.domain()

with open('/tmp/on_domain_change.log', 'w+') as f:
    f.write(domain)
