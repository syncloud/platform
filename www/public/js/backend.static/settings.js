backend.device_data = {
      "device_domain": "test.syncloud.it",
      "success": true
    };

backend.access_data = {
      "data": {
        "external_access": true,
        "protocol": "https"
      },
      "success": true
    };

backend.versions_data = {
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
    };

backend.disks_data = {
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
    };

backend.boot_disk_data = {
      "data": {
          "device": "/dev/mmcblk0p2",
          "size": "2G",
          "extendable": true
        },
      "success": true
    };

backend.device_url = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() { on_complete(that.device_data); }, 2000);
    };

backend.send_logs = function(include_support, on_always, on_error) {
        setTimeout(on_always, 2000);
    };

backend.reactivate = function(on_complete, on_error) {
        on_complete({ activate_url: "../internal/activate.html", success: true});
    };

backend.get_versions = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() { 
            on_complete(that.versions_data);
        }, 2000);
    };

backend.check_versions = function(on_always, on_error) {
        setTimeout(on_always, 2000);
    };

backend.platform_upgrade = function(on_complete, on_error) {
        setTimeout(on_complete({success: true}), 2000);
    },

backend.boot_extend = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() {
            that.boot_disk_data.data.extendable = false;
            that.boot_disk_data.data.size = '16G';
            on_complete({success: true});
        }, 2000);
    };

backend.sam_upgrade = function(on_complete, on_error) {
        setTimeout(function() { on_complete({success: true}) }, 2000);
    };

backend.update_disks = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() { on_complete(that.disks_data); }, 2000);
    };

backend.update_boot_disk = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() { on_complete(that.boot_disk_data); }, 2000);
    };

backend.disk_action = function(disk_device, is_activate, on_complete, on_error) {
        var that = this;
        setTimeout(function() { on_complete(that.disks_data); }, 2000);
    };