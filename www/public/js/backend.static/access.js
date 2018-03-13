backend.access_data = {
        error_toggle: false,
        "data": {
            "external_access": true,
            "upnp_available": false,
            "upnp_enabled": true,
            "upnp_message": "Your router does not have port mapping feature enabled at the moment",
            "public_ip": null
        },
        "success": true
    };

backend.network_interfaces_data = {
        "data": {
            "interfaces": [
                {
                    "ipv4": [
                        {
                            "addr": "172.17.0.2",
                            "broadcast": "172.17.0.2",
                            "netmask": "255.255.0.0"
                        }
                    ],
                    "ipv6": [
                        {
                            "addr": "fe80::42:acff:fe11:2%eth0",
                            "netmask": "ffff:ffff:ffff:ffff::"
                        }
                    ],
                    "name": "eth0"
                }
            ]
        },
        "success": true
    };
backend.port_mappings_data = {
        "port_mappings": [
             {
                 "local_port": 80,
                 "external_port": 80                 
             },
             {
                 "local_port": 443,
                 "external_port": 10001                     
             }
                    
        ],
        "success": true
    };

backend.check_access = function (on_complete, on_error) {
        var that = this;
        setTimeout(function () { on_complete(that.access_data); }, 2000);
    };

backend.set_access = function (upnp_enabled,
                               external_access,
                               public_ip,
                               public_port,
                               on_complete,
                               on_error) {
        var that = this;
        setTimeout(function () {
            if (that.access_data.error_toggle) {
                that.access_data.data.external_access = external_access;
                that.access_data.data.upnp_enabled = upnp_enabled;
                that.access_data.data.public_ip = public_ip;
                that.access_data.data.public_port = public_port;
                on_complete({success: true});
            } else {
                var xhr = {
                    status: 200,
                    responseJSON: {
                        message: "error"
                    }
                };
                on_error(xhr, {}, {});
            }
            that.access_data.error_toggle = ! that.access_data.error_toggle;
        }, 2000);
    };


backend.network_interfaces = function (on_complete, on_error) {
        var that = this;
        setTimeout(function () { on_complete(that.network_interfaces_data); }, 2000);
    };
    
backend.port_mappings = function (on_complete, on_error) {
        var that = this;
        setTimeout(function () { on_complete(that.port_mappings_data); }, 2000);
    };