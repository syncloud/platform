from syncloud_platform.gaplib import gen

from os import remove
from test.gaplib.helpers import temp_file, read_text
from os.path import isfile

import pytest

def test_empty_template():
    empty_template_filename = temp_file();
    target_filename = temp_file();
    remove(target_filename)
    gen.generate_file_jinja(empty_template_filename, target_filename, {})
    assert isfile(target_filename)


def test_simple_substitution():
    template_text = '{{ variable }}'
    empty_template_filename = temp_file(text=template_text);
    target_filename = temp_file();
    remove(target_filename)
    gen.generate_file_jinja(empty_template_filename, target_filename, {'variable': 'blah-blah'})
    assert isfile(target_filename)
    target_text = read_text(target_filename)
    assert target_text == 'blah-blah'


def test_unknown_variable():
    template_text = '{{ variable }}'
    empty_template_filename = temp_file(text=template_text);
    target_filename = temp_file();
    remove(target_filename)
    gen.generate_file_jinja(empty_template_filename, target_filename, {})
    assert isfile(target_filename)
    target_text = read_text(target_filename)
    assert target_text == template_text


def test_custom_variable_tags():
    template_text = '<% variable > brackets {{ variable }} does not mean anything'
    empty_template_filename = temp_file(text=template_text);
    target_filename = temp_file();
    remove(target_filename)
    gen.generate_file_jinja(empty_template_filename, target_filename, {'variable': 'blah-blah'}, variable_tags=('<%', '>'))
    assert isfile(target_filename)
    target_text = read_text(target_filename)
    assert target_text == 'blah-blah brackets {{ variable }} does not mean anything'
