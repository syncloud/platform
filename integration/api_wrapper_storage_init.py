import sys, os
from syncloudlib.application.storage import get_storage_dir

print os.path.realpath(get_storage_dir(sys.argv[1]))