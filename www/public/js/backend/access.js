var backend = {

    check_access: function(on_complete, on_error) {
        $.get('/rest/access/access').done(on_complete).fail(on_error);
    },

    set_access: function(
        upnp_enabled,
        external_access,
        is_https,
        public_ip,
        public_port,
        on_complete,
        on_error) {
        
        $.get('/rest/access/set_access', {
            upnp_enabled: upnp_enabled,
            external_access: external_access,
            is_https: is_https,
            public_ip: public_ip,
            public_port: public_port
        }).done(on_complete).fail(on_error);
    },

    network_interfaces: function(on_complete, on_error) {
        $.get('/rest/access/network_interfaces').done(on_complete).fail(on_error);
    }

};