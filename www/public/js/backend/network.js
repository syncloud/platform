backend.check_access = function(on_complete, on_error) {
        $.get('/rest/access/access').done(on_complete).fail(on_error);
    };

backend.set_access = function(data, on_complete, on_error) {
        $.get('/rest/access/set_access', data).done(on_complete).fail(on_error);
    };
    
backend.network_interfaces = function(on_complete, on_error) {
        $.get('/rest/access/network_interfaces').done(on_complete).fail(on_error);
    };
    
backend.port_mappings = function(on_complete, on_error) {
        $.get('/rest/access/port_mappings').done(on_complete).fail(on_error);
    };
