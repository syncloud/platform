import requests_unixsocket

socket_file = '/var/snap/platform/current/backend.sock'
socket = 'http+unix://{0}'.format(socket_file.replace('/', '%2F'))


def backend_request(method, url, data):
    session = requests_unixsocket.Session()
    return session.request(method, '{0}{1}'.format(socket, url), json=data)
