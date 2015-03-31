from subprocess import check_output


def local_ip():
    local_ip = check_output(["hostname", "-I"]).split(" ")[0]
    if not local_ip:
        raise(Exception("Can't get local ip address"))
    return local_ip


