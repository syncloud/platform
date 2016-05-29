import os


def restart(delay_sec=0):
    command = '( sleep {delay_sec}; systemctl reboot -i ) &'.format(delay_sec=delay_sec)
    os.system(command)


def shutdown(delay_sec=0):
    command = '( sleep {delay_sec}; systemctl poweroff -i ) &'.format(delay_sec=delay_sec)
    os.system(command)