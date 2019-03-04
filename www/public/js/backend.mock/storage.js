backend.disks_data = {
      "disks": [
        {
          "name": "My Passport 0837",
          "device": "/dev/sdb",
          "active": true,
          "size": "931.5G",
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
          "device": "/dev/sdc",
          "active": false,
          "size": "931.5G",
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
        },
        {
          "name": "Blank Disk",
          "device": "/dev/sdb",
          "size": "100 TB",
          "partitions": []
        }
      ],
      "success": true
    };

backend.disks_data_error = {
      "message": "error",
      "success": false
    };

backend.boot_disk_data = {
      "data": {
          "device": "/dev/mmcblk0p2",
          "size": "2G",
          "extendable": true
        },
      "success": true
    };

backend.boot_extend = function(on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() {
            that.boot_disk_data.data.extendable = false;
            that.boot_disk_data.data.size = '16G';
            on_complete({success: true});
        }, 2000);
    };

backend.update_disks = function(on_complete, on_error) {
        var that = this;
        setTimeout(function() { on_complete(that.disks_data); }, 500);
    };

backend.update_boot_disk = function(on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() { on_complete(that.boot_disk_data); }, 2000);
    };

backend.disk_action_success = true;

backend.disk_action = function(disk_device, is_activate, on_complete, on_always, on_error) {
        var that = this;
        if (backend.disk_action_success) {
            backend.test_timeout(function() { on_complete(that.disks_data); on_always(); }, 2000);
        } else {
            backend.test_timeout(function() { on_complete(that.disks_data_error); on_always(); }, 2000);
        }
    };
    
backend.disk_format = function(disk_device, on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() {
            on_complete({success: true});
        }, 2000);
    };