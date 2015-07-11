from subprocess import check_output

class Cmd:
    def run(self, command_line):
        return check_output(command_line, shell=True)
