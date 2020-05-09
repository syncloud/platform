import sys, os
from syncloudlib.application.storage import init_storage

print(os.path.realpath(init_storage(sys.argv[1], sys.argv[2])))