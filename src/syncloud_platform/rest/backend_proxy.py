import requests_unixsocket

socket_file = '/var/snap/platform/common/backend.sock'
socket = 'http+unix://{0}'.format(socket_file.replace('/', '%2F'))

def backend_get(url):
    session = requests_unixsocket.Session()
    return session.get('{0}{1}'.format(socket, url))
