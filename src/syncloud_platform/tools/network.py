from subprocess import check_output


class Network:

    def __init__(self):
        pass

    def local_ip(self):
        local_ip = check_output(["hostname", "-I"]).split(" ")[0]
        if not local_ip:
            raise(Exception("Can't get local ip address"))
        return local_ip
