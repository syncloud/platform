#!/opt/app/platform/python/bin/python
from syncloud_platform.application import api

app = api.get_app_setup('platform')
domain = api.device_domain_name()

with open('/tmp/on_domain_change.log', 'w+') as f:
    f.write(domain)
