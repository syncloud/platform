var backend = {
    device_data: {
      "device_domain": "test.syncloud.it",
      "success": true
    },

    access_data: {
      "data": {
        "external_access": true,
        "protocol": "https"
      },
      "success": true
    },

    versions_data: {
      "data": [
        {
          "app": {
            "id": "platform",
            "name": "Platform",
            "required": true,
            "ui": false,
            "url": "http://platform.odroid-c2.syncloud.it"
          },
          "current_version": "880",
          "installed_version": "876"
        },
        {
          "app": {
            "id": "sam",
            "name": "Syncloud App Manager",
            "required": true,
            "ui": false,
            "url": "http://sam.odroid-c2.syncloud.it"
          },
          "current_version": "78",
          "installed_version": "75"
        }
      ],
      "success": true
    },

    disks_data: {
      "disks": [
        {
          "name": "My Passport 0837",
          "partitions": [
            {
              "active": true,
              "device": "/dev/sdb1",
              "fs_type": "ntfs",
              "mount_point": "/opt/disk/external",
              "mountable": true,
              "size": "931.5G"
            }
          ]
        },
        {
          "name": "My Passport 0990",
          "partitions": [
            {
              "active": false,
              "device": "/dev/sdc1",
              "fs_type": "ntfs",
              "mount_point": "",
              "mountable": true,
              "size": "931.5G"
            }
          ]
        }
      ],
      "success": true
    },

    boot_disk_data: {
      "data": {
          "device": "/dev/mmcblk0p2",
          "size": "2G",
          "extendable": true
        },
      "success": true
    },

    device_url: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.device_data);
        }, 2000);
    },

    send_logs: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    },

    reactivate: function() {
        window.location.href = "activate.html";
    },

    check_access: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.access_data);
        }, 2000);
    },

    external_access: function(parameters) {
        var that = this;
        setTimeout(function() {
            that.access_data.data.external_access = parameters.state;
            if (!that.access_data.data.external_access) {
                that.access_data.data.protocol = "http";
            }
            success_callbacks(parameters);
        }, 2000);
    },

    protocol: function(parameters) {
        var that = this;
        setTimeout(function() {
            that.access_data.data.protocol = parameters.new_protocol;
            success_callbacks(parameters);
        }, 2000);
    },

    get_versions: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.versions_data);
        }, 2000);
    },

    check_versions: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    },

    platform_upgrade: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    },

    boot_extend: function(parameters) {
        var that = this;
        setTimeout(function() {
            that.boot_disk_data.data.extendable = parameters.extendable;
            that.boot_disk_data.data.size = '16G';
            success_callbacks(parameters);
        }, 2000);
    },

    sam_upgrade: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    },

    update_disks: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.disks_data);
        }, 2000);
    },

    update_boot_disk: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.boot_disk_data);
        }, 2000);
    },

    disk_action: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.disks_data);
        }, 2000);
    }

};