import os
from os import makedirs
from os.path import join, dirname, relpath, isdir, split
from string import Template
import jinja2

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

def generate_file_jinja(from_path, to_path, variables, variable_tags=('{{', '}}')):
    from_path_dir, from_path_filename = split(from_path)
    loader = jinja2.FileSystemLoader(searchpath=from_path_dir)
    variable_start_tag, variable_end_tag = variable_tags
    env_parameters = dict(
        loader=loader,
        trim_blocks=True,
        lstrip_blocks=True,
        undefined=jinja2.StrictUndefined,
        variable_start_string=variable_start_tag,
        variable_end_string=variable_end_tag
    )
    environment = jinja2.Environment(**env_parameters)
    template = environment.get_template(from_path_filename)
    output = template.render(variables)
    to_path_dir = dirname(to_path)
    if not isdir(to_path_dir):
        makedirs(to_path_dir)
    with open(to_path, 'wb+') as fh:
        fh.write(output)

def generate_files(from_dir, to_dir, variables, variable_tags=('{{', '}}')):
    for dir_name, subdirs, files in os.walk(from_dir):
        for filename in files:
            from_path = join(dir_name, filename)
            from_rel_path = relpath(from_path, from_dir)
            to_path = join(to_dir, from_rel_path)
            generate_file_jinja(from_path, to_path, variables, variable_tags)
