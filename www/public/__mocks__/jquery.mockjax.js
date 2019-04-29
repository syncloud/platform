var jquery = require('jquery');
var mockjax = require('jquery-mockjax')(jquery, window);

export const State = {
    available_apps_success: true,
    disk_action_success: true
}

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


const access_data = {
        error_toggle: false,
        "data": {
            "external_access": true,
            "upnp_available": false,
            "upnp_enabled": true,
            "upnp_message": "not used",
            "public_ip": "111.111.111.111"
        },
        "success": true
    };

const network_interfaces_data = {
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

const port_mappings_data = {
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

mockjax({
    url: '/rest/access/access',
    dataType: "json",
    responseText: access_data
});

function set_access(settings) {
    if (!access_data.error_toggle) {
        access_data.data.external_access = settings.data.external_access;
        access_data.data.upnp_enabled = settings.data.upnp_enabled;
        if (settings.data.ip_autodetect) {
            if (access_data.data.hasOwnProperty('public_ip')) {
                delete access_data.data.public_ip;
            }
        } else {
            access_data.data.public_ip = settings.data.public_ip;
        }
        if (settings.data.upnp_enabled) {
            port_mappings_data.port_mappings[0].external_port = 81;
            port_mappings_data.port_mappings[1].external_port = 444;
        } else {
            port_mappings_data.port_mappings[0].external_port = settings.data.certificate_port;
            port_mappings_data.port_mappings[1].external_port = settings.data.access_port;
        }
        this.responseText = {
            "success": true
        };
    } else {
        this.responseText = {
            "success": false,
            "message": "error"
        };
    }
    access_data.error_toggle = ! access_data.error_toggle;
};

mockjax({
    url: '/rest/access/set_access',
    dataType: "json",
    response: set_access
});

mockjax({
    url: '/rest/access/network_interfaces',
    dataType: "json",
    responseText: network_interfaces_data
});

mockjax({
    url: '/rest/access/port_mappings',
    dataType: "json",
    responseText: port_mappings_data
});

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

function disk_response(settings) {
      if (State.disk_action_success) {
        this.responseText = disks_data;
      } else {
        this.responseText = disks_data_error;
      }
    }

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
 
mockjax({ url:'/rest/settings/boot_extend_status', dataType: "json",
    responseText: {success: true, is_running: false}
});

mockjax({ url:'/rest/settings/disk_format_status',dataType: "json",
    responseText: {success: true, is_running: false}
});

function sam_status(settings) {
    console.log('sam status mock');
    this.responseText = {success: true, is_running: false};
}

mockjax({ 
    url:'/rest/settings/sam_status',
    dataType: "json",
    response: sam_status
});

export const versions_data = {
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

mockjax({ url:'/rest/settings/versions', dataType: "json",
    responseText: versions_data
});

mockjax({ url:'/rest/check', dataType: "json",
    responseText: {}
});

mockjax({ url:'/rest/upgrade', dataType: "json",
    responseText: {success: true}
});

const appcenter_data = {
      "apps": [
        {
          "id": "owncloud",
          "name": "ownCloud",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "diaspora",
          "name": "Diaspora",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "mail",
          "name": "Mail",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "talk",
          "name": "Talk",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "files",
          "name": "Files Browser",
          "icon": "appsimages/penguin.png"
        }
      ]
    };

const appcenter_data_error = {
      "message": "error",
      "success": false
    };

function available_apps_response(settings) {
        if (State.available_apps_success) {
            this.responseText = appcenter_data;
        } else {
            this.responseText = appcenter_data_error;
        }
    }

mockjax({ url:'/rest/available_apps', dataType: "json",
    response: available_apps_response
});

mockjax({ url:'/rest/app', dataType: "json",
    responseText: {
        info: {
            "app": {
                "id": "wordpress",
                "name": "Wordpress",
                "required": true,
                "ui": false,
                "url": "http://platform.odroid-c2.syncloud.it",
                "icon": "appsimages/penguin.png"
            },
            "current_version": "880",
            "installed_version": "876"
        }
    }
});

const device_data = {
  "device_url": "http://test.syncloud.it",
  "success": true
};

mockjax({ url:'/rest/settings/device_url', dataType: "json",
    responseText: device_data
});

var backups = [
            {"path":"/data/platform/backup","file":"files-2019-0515-123506.tar.gz"},
            {"path":"/data/platform/backup","file":"nextcloud-2019-0515-123506.tar.gz"},
            {"path":"/data/platform/backup","file":"files-2019-0415-123506.tar.gz"}
        ];

function backupList(settings) {
    //alert(backups.length);
    this.responseText = {
        "success":true,
        "data": backups
    };
}

mockjax({ 
    url:'/rest/backup/list',
    dataType: "json",
    response: backupList
});

function backupRemove(settings) {
    //alert("remove "+settings.data.file);
    backups = backups.filter(v => v.file != settings.data.file);
    //alert(backups);
}

mockjax({ 
    type: "post",
    url:'/rest/backup/remove', dataType: "json",
    response: backupRemove
});

function backupRestore(settings) {
    //alert("restore "+settings.data.file);
}

mockjax({ 
    type: "post",
    url:'/rest/backup/restore', dataType: "json",
    response: backupRestore
});

function backupCreate(settings) {
    //alert("restore "+settings.data.file);
}

mockjax({ 
    type: "post",
    url:'/rest/backup/create', dataType: "json",
    response: backupCreate
});
