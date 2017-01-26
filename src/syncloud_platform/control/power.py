import os


def restart():
    command = 'shutdown -r now'
    os.system(command)


def shutdown():
    command = 'shutdown now'
    os.system(command)