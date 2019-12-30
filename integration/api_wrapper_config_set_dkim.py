import sys
from syncloudlib.application.config import get_dkim_key, set_dkim_key

print set_dkim_key(sys.argv[1])
