import pytest
from syncloud_platform.insider.upnpc import check_error


def test_check_error_throws():
    with pytest.raises(Exception) as excinfo:
        check_error("""
upnpc : miniupnpc library test client. (c) 2006-2010 Thomas Bernard
Go to http://miniupnp.free.fr/ or http://miniupnp.tuxfamily.org/
for more information.
List of UPNP devices found on the network :
 desc: http://192.168.1.254:2555/upnp/UPnP_BThomeHub5A_18622c3d47be/desc.xml
 st: urn:schemas-upnp-org:device:InternetGatewayDevice:1

Found valid IGD : http://192.168.1.254:2555/upnp/UPnP_BThomeHub5A_18622c3d47be_ptm0/WANPPPConn1.ctl
Local LAN ip address : 192.168.1.73
ExternalIPAddress = 86.153.151.12
AddPortMapping(1023, 80, 192.168.1.73) failed with code 501 (Action Failed)
GetSpecificPortMappingEntry() failed with code -1 (Miniupnpc Unknown Error)""")

    assert 'Unable to add mapping: AddPortMapping(1023, 80, 192.168.1.73) failed with code 501 (Action Failed)' in str(excinfo.value)


def test_check_error_returns():
    input = """
upnpc : miniupnpc library test client. (c) 2006-2010 Thomas Bernard
Go to http://miniupnp.free.fr/ or http://miniupnp.tuxfamily.org/
for more information.
List of UPNP devices found on the network :
 desc: http://192.168.1.254:2555/upnp/UPnP_BThomeHub5A_18622c3d47be/desc.xml
 st: urn:schemas-upnp-org:device:InternetGatewayDevice:1

Found valid IGD : http://192.168.1.254:2555/upnp/UPnP_BThomeHub5A_18622c3d47be_ptm0/WANPPPConn1.ctl
Local LAN ip address : 192.168.1.73
ExternalIPAddress = 86.153.151.12
InternalIP:Port = 192.168.1.73:80
external 86.153.151.12:1031 TCP is redirected to internal 192.168.1.73:80"""

    assert check_error(input) == input