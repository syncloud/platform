import os
from os.path import join, dirname, relpath
from string import Template

import fs

def generate_file(from_path, to_path, variables):
    from_file = open(from_path, 'r')
    from_text = from_file.read()
    from_file.close()
    t = Template(from_text)
    to_text = t.substitute(variables)
    fs.makepath(dirname(to_path))
    to_file = open(to_path, 'w+')
    to_file.write(to_text)
    to_file.close()

def generate_files(from_dir, to_dir, variables):
    for dir_name, subdirs, files in os.walk(from_dir):
        for filename in files:
            from_path = join(dir_name, filename)
            from_rel_path = relpath(from_path, from_dir)
            to_path = join(to_dir, from_rel_path)
            generate_file(from_path, to_path, variables)
