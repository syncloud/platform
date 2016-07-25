var device_data = {
  "device_domain": "test.syncloud.it",
  "success": true
}

var access_data = {
  "data": {
    "external_access": false,
    "protocol": "http"
  },
  "success": true
};

var disks_data = {
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

var versions_data = {
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

function backend_device_url(on_complete) {
    setTimeout(function() {
        display_device_url(device_data);
        on_complete();
    }, 2000);
}

function backend_send_logs(on_complete) {
    setTimeout(function() {
        on_complete();
    }, 2000);
}

function backend_check_versions(on_complete) {
    setTimeout(function() {
        display_versions(versions_data);
        on_complete();
    }, 2000);
}

function backend_check_access(on_complete) {
    setTimeout(function() {
        display_access(access_data.data);
        on_complete();
    }, 2000);
}

function backend_reactivate() {
    window.location.href = "activate.html";
}

function backend_external_access(state, on_complete) {
    setTimeout(function() {
        access_data.data.external_access = state;
        display_access(access_data.data);
        on_complete();
    }, 2000);
}

function backend_protocol(new_protocol, on_complete) {
    setTimeout(function() {
        access_data.data.protocol = new_protocol;
        display_access(access_data.data);
        on_complete();
    }, 2000);
}

function backend_platform_upgrade(on_complete) {
    setTimeout(function() {
        display_versions(versions_data);
        on_complete();
    }, 2000);
}

function backend_sam_upgrade(on_complete) {
    setTimeout(function() {
        display_versions(versions_data);
        on_complete();
    }, 2000);
}

function backend_update_disks(on_complete) {
    setTimeout(function() {
        display_disks(disks_data);
        on_complete();
    }, 2000);
}

function backend_disk_action(disk_device, is_activate, on_complete) {
    setTimeout(function() {
        var disks = disks_data.disks;
        for (i=0; i < disks.length; i++) {
            var disk = disks[i];
            partitions = disk.partitions;
            for (j=0; j < partitions.length; j++) {
                var partition = partitions[j];
                if (partition.device == disk_device) {
                    partition.active = is_activate;
                } else if (is_activate) {
                    partition.active = false;
                }
            }
        }
        display_disks(disks_data);
    }, 2000);
}