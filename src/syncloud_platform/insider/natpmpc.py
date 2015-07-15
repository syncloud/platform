import NATPMP

class NatPmpPortMapper:
    def external_ip(self):
        return NATPMP.get_public_address()

    def add_mapping(self, local_port):
        response = NATPMP.map_port(NATPMP.NATPMP_PROTOCOL_TCP, local_port, local_port)
        return response.public_port

    def remove_mapping(self, local_port, external_port):
        NATPMP.map_port(NATPMP.NATPMP_PROTOCOL_TCP, local_port, external_port, lifetime=0)