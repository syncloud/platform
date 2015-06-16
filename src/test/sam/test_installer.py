from os.path import dirname
from tempfile import mkstemp
from syncloud.sam.installer import fix_locale_gen

test_dir = dirname(__file__)
original_config_file = test_dir + '/data/locale.gen.original'
expected_config_file = test_dir + '/data/locale.gen.fixed'


def test_lang_to_replace():

    from_fh, actual_config_file = mkstemp()
    with open(actual_config_file, 'w') as f:
        f.write(open(original_config_file).read())

    fix_locale_gen('en_GB.UTF-8', actual_config_file)

    assert open(actual_config_file).read() == open(expected_config_file).read()
