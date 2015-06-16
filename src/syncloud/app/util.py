from string import Template


def transform_file(from_filename, to_filename, mapping):
    with open(from_filename, 'r') as from_f:
        template = Template(from_f.read())
        runtime = template.substitute(mapping)
        with open(to_filename, 'w') as to_f:
            to_f.write(runtime)
    return to_filename
