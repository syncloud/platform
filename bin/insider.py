from syncloud_platform.injector import get_injector
import argparse
from syncloud_app import main


def create_args_parser():
    parser = argparse.ArgumentParser(description='Syncloud insider maps port on router and creates DNS records')
    parser.add_argument('--debug', action='store_true')

    subparsers = parser.add_subparsers(help='available commands', dest='action')
    subparsers.add_parser('sync_all', help="sync port mappings and dns records")
    return parser


if __name__ == '__main__':
    parser = create_args_parser()
    args = parser.parse_args()

    device = get_injector(debug=args.debug).device

    main.execute(device, args)
