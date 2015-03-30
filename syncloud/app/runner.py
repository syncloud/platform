import subprocess
import select
from logging import DEBUG, ERROR


#https://gist.github.com/bgreenlee/1402841
def call(popenargs, logger, stdout_log_level=DEBUG, stderr_log_level=ERROR, **kwargs):
    """
    Variant of subprocess.call that accepts a logger instead of stdout/stderr,
    and logs stdout messages via logger.debug and stderr messages via
    logger.error.
    """
    logger.info(popenargs)
    child = subprocess.Popen(popenargs, stdout=subprocess.PIPE, stderr=subprocess.PIPE, **kwargs)

    log_level = {child.stdout: stdout_log_level,
                 child.stderr: stderr_log_level}

    def check_io():
        ready_to_read = select.select([child.stdout, child.stderr], [], [], 1000)[0]
        for io in ready_to_read:
            line = io.readline()
            if not (io == child.stderr and not line):
                line_strip = line.strip()
                if line_strip:
                    logger.log(log_level[io], line_strip)

    # keep checking stdout/stderr until the child exits
    while child.poll() is None:
        check_io()

    check_io()  # check again to catch anything after the process exits

    return child.wait()