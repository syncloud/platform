from os.path import dirname, abspath
import sys


def run_script(script_filename, add_location_to_sys_path=False):
    if add_location_to_sys_path:
        script_folder = abspath(dirname(script_filename))
        sys.path.append(script_folder)
    g = globals().copy()
    g['__file__'] = script_filename
    execfile(script_filename, g)
