var backend = {

    access_data: {
        "data": {
            "external_access": true,
            "protocol": "https",
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

    check_access: function (parameters) {
        var that = this;
        setTimeout(function () {
            success_callbacks(parameters, that.access_data);
        }, 2000);
    },

    save_access: function (parameters) {
        var that = this;
        setTimeout(function () {
            that.access_data.data.external_access = parameters.access;
            that.access_data.data.upnp_enabled = parameters.upnp;
            that.access_data.data.public_ip = parameters.public_ip;
            that.access_data.data.public_port = parameters.public_port;
            if (!that.access_data.data.external_access) {
                that.access_data.data.protocol = "http";
            }
            success_callbacks(parameters);
        }, 2000);
    },

    protocol: function (parameters) {
        var that = this;
        setTimeout(function () {
            that.access_data.data.protocol = parameters.new_protocol;
            success_callbacks(parameters);
        }, 2000);
    },

    network_interfaces: function (parameters) {
        var that = this;
        setTimeout(function () {
            success_callbacks(parameters, that.network_interfaces_data);
        }, 2000);
    }
};