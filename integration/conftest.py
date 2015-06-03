import pytest


def pytest_addoption(parser):
    parser.addoption("--email", action="store")
    parser.addoption("--password", action="store")


@pytest.fixture(scope="session")
def auth(request):
    config = request.config
    return config.getoption("--email"), config.getoption("--password")