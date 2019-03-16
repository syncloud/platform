backend.check_access = function(on_complete, on_error) {
        $.get('/rest/access/access').done(on_complete).fail(on_error);
    };

backend.set_access = function(
        upnp_enabled,
        external_access,
        ip_autodetect,
        public_ip,
        certificate_port,
        access_port,
        on_complete,
        on_error) {

        var request_data = {
           upnp_enabled: upnp_enabled,
           external_access: external_access,
           certificate_port: certificate_port,
           access_port: access_port
        };

        if (!ip_autodetect) {
            request_data.public_ip = public_ip;
        }

        $.get('/rest/access/set_access', request_data).done((data) => { alert(data); on_complete(data); }).fail(on_error);
    };
    
backend.network_interfaces = function(on_complete, on_error) {
        $.get('/rest/access/network_interfaces').done(on_complete).fail(on_error);
    };
    
backend.port_mappings = function(on_complete, on_error) {
        $.get('/rest/access/port_mappings').done(on_complete).fail(on_error);
    };
