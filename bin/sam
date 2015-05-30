#!/usr/bin/env python
import httplib
import logging
from os.path import join, dirname
import sys
import argparse
from syncloud.app import main
from syncloud.app import logger

sys.path.append(join(dirname(__file__), '..'))

from syncloud.sam.manager import get_sam


def get_arg_parser():
    parser = argparse.ArgumentParser(description='Syncloud application manager')
    parser.add_argument('--debug', action='store_true')

    subparsers = parser.add_subparsers(help='available commands', dest='action')

    subparsers.add_parser('list', help="list apps")

    sub = subparsers.add_parser('install', help="install application")
    sub.add_argument('app_id', help="application id")

    sub = subparsers.add_parser('remove', help="remove application")
    sub.add_argument('app_id', help="application id")

    sub = subparsers.add_parser('verify', help="verify application")
    sub.add_argument('app_id', help="application id")

    sub = subparsers.add_parser('latest', help="latest version of application")
    sub.add_argument('app_id', help="application id")

    sub = subparsers.add_parser('update', help="update apps repository")
    sub.add_argument('--release', default=None, dest='release')

    sub = subparsers.add_parser('upgrade_all', help="upgrade apps and install required apps")

    return parser

if __name__ == '__main__':
    arg_parser = get_arg_parser()
    args = arg_parser.parse_args()

    console = True if args.debug else False
    level = logging.DEBUG if args.debug else logging.INFO
    logger.init(level, console, '/var/log/sam.log')

    # if args.debug:
    #     httplib.HTTPConnection.debuglevel = 1

    sam = get_sam()

    main.execute(sam, args)