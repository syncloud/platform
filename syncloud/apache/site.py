import glob
import os
import shutil
from syncloud.app.logger import get_logger
from env import http_include_dir_full, https_include_dir_full
from model import Site
from ports import Ports


def get_port(https):
    ports = Ports()
    if https:
        return ports.https
    return ports.http


def add(name, config_file, https):
    log = get_logger("apache.add")
    if not os.path.exists(config_file):
        raise Exception("{} does not exist".format(config_file))

    remove(name, https)

    name_conf_full = full_conf(name, https)
    log.info('copy file {0} to {1}'.format(config_file, name_conf_full))
    shutil.copyfile(config_file, name_conf_full)

    return get_port(https)


def remove(name, https):
    name_conf_full = full_conf(name, https)
    if os.path.exists(name_conf_full):
        os.remove(name_conf_full)


def full_conf(name, https):

    if https:
        include_dir_full = https_include_dir_full
    else:
        include_dir_full = http_include_dir_full

    name_conf = "{}.conf".format(name)
    return os.path.join(include_dir_full, name_conf)


def _merge_sites(sites, include_dir, protocol):
    for site in _read_sites(include_dir):
        if not site.name in sites:
            sites[site.name] = site
        sites[site.name].protocols.append(protocol)


def _read_sites(include_dir):
    return [_create_site(conf_file) for conf_file in _read_confs(include_dir)]


def _read_confs(conf_dir):
    return glob.glob("{0}/*.conf".format(conf_dir))


def _create_site(conf_file):
    site_path, conf_file = os.path.split(conf_file)
    site_name = os.path.splitext(conf_file)[0]
    return Site(site_name)