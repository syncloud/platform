from subprocess import check_output, CalledProcessError
from os.path import join
from syncloud.sam.manager import default_config_path, config_file


def test_run():
    print open(join(default_config_path, config_file), 'r').read()

    try:
        print check_output(['syncloud-cli', '-h'])
    except CalledProcessError, e:
        print(e.output)
        raise e