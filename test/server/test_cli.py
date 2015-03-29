from subprocess import check_output, CalledProcessError


def test_run():
    try:
        print check_output(['syncloud-cli', '-h'])
    except CalledProcessError, e:
        print(e.output)
        raise e