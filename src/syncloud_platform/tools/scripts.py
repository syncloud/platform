def run_script(script_filename):
    g = globals().copy()
    g['__file__'] = script_filename
    execfile(script_filename, g)
