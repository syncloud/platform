import sys
from syncloudlib.application.config import set_dkim_key

print(set_dkim_key(sys.argv[1]))
