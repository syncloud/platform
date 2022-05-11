from subprocess import check_output


def parted(device):
    return check_output('parted {0} unit % print free --script --machine'.format(device).split(" ")).decode()
