var jquery = require('jquery');
var mockjax = require('jquery-mockjax')(jquery, window);

const disks_data = {
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

const disks_data_error = {
      "message": "error",
      "success": false
    };

const boot_disk_data = {
      "data": {
          "device": "/dev/mmcblk0p2",
          "size": "2G",
          "extendable": true
        },
      "success": true
    };

var disk_action_success = true;

mockjax({
    url: '/rest/settings/disks',
    dataType: "json",
    responseText: disks_data
});

mockjax({
    url: '/rest/settings/boot_disk',
    dataType: "json",
    responseText: boot_disk_data
});

const disk_response = (settings) => {
      if (disk_action_success) {
        this.responseText = disks_data;
      } else {
        this.responseText = disks_data_error;
      }
    }

const apps_data = {
    "apps": [
        {
            "id": "wordpress",
            "name": "WordPress",
            "icon": "appsimages/penguin.png",
            "url": "http://owncloud.odroid-c2.syncloud.it"
        }
    ]
};

mockjax({
    url: '/rest/installed_apps',
    dataType: "json",
    responseText: apps_data
});

mockjax({
    url: '/rest/settings/disk_activate',
    dataType: "json",
    response: disk_response
});
    
mockjax({
    url: '/rest/settings/disk_deactivate',
    dataType: "json",
    response: disk_response
});

mockjax({
    url:'/rest/settings/boot_extend',
    dataType: "json",
    responseText: {success: true}
});
    
mockjax({
    url:'/rest/storage/disk_format',
    dataType: "json",
    responseText: {success: true}
});
 
mockjax({
    url:'/rest/settings/boot_extend_status',
    dataType: "json",
    responseText: {success: true, is_running: false}
});

mockjax({
    url:'/rest/settings/disk_format_status',
    dataType: "json",
    responseText: {success: true, is_running: false}
});
