from os import makedirs
from os.path import dirname, isdir, split

import jinja2


def generate_file_jinja(from_path, to_path, variables, variable_tags=('{{', '}}')):
    from_path_dir, from_path_filename = split(from_path)
    loader = jinja2.FileSystemLoader(searchpath=from_path_dir)
    variable_start_tag, variable_end_tag = variable_tags
    env_parameters = dict(
        loader=loader,
        # some files like udev rules want empty lines at the end
        # trim_blocks=True,
        # lstrip_blocks=True,
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
        fh.write(output.encode("UTF-8"))

