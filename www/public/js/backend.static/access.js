var backend = {

    access_data: {
        "data": {
            "external_access": true,
            "is_https": false,
            "upnp_available": false,
            "upnp_enabled": true,
            "upnp_message": "Your router does not have port mapping feature enabled at the moment",
            "public_ip": null
        },
        "success": true
    },

    network_interfaces_data: {
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
    },

    check_access: function (on_complete, on_error) {
        var that = this;
        setTimeout(function () { on_complete(that.access_data); }, 2000);
    },

    set_access: function (upnp_enabled,
                          external_access,
                          is_https,
                          public_ip,
                          public_port,
                          on_complete,
                          on_error) {
        var that = this;
        setTimeout(function () {
            that.access_data.data.external_access = external_access;
            that.access_data.data.upnp_enabled = upnp_enabled;
            that.access_data.data.public_ip = public_ip;
            that.access_data.data.public_port = public_port;
            that.access_data.data.is_https = is_https;
            on_complete({success: true});
        }, 2000);
    },


    network_interfaces: function (on_complete, on_error) {
        var that = this;
        setTimeout(function () { on_complete(that.network_interfaces_data); }, 2000);
    }
};