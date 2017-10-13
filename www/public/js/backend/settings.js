backend.device_url = function(on_complete, on_error) {
        $.get('/rest/settings/device_url').done(on_complete).fail(on_error);
    };

backend.send_logs = function(include_support, on_always, on_error) {
        $.get('/rest/send_log',
          { include_support: include_support }
        ).always(on_always).fail(on_error);
    };

backend.reactivate = function(on_complete, on_error) {
        $.get('/rest/settings/activate_url').done(on_complete).fail(on_error);
    };

backend.get_versions = function(on_complete, on_error) {
        $.get('/rest/settings/versions').done(on_complete).fail(on_error);
    };

backend.check_versions = function(on_always, on_error) {
        $.get('/rest/check').always(on_always).fail(on_error);
    };

backend.platform_upgrade = function(on_complete, on_error) {
        $.get('/rest/settings/system_upgrade').done(on_complete).fail(on_error);
    };

backend.sam_upgrade = function(on_complete, on_error) {
        $.get('/rest/settings/sam_upgrade').done(on_complete).fail(on_error);
    };

backend.update_disks = function(on_complete, on_error) {
        $.get('/rest/settings/disks').done(on_complete).fail(on_error);
    };

backend.update_boot_disk = function(on_complete, on_error) {
        $.get('/rest/settings/boot_disk').done(on_complete).fail(on_error);
    };

backend.disk_action = function(disk_device, is_activate, on_complete, on_always, on_error) {
        var mode = is_activate ? "disk_activate" : "disk_deactivate";
        $.get('/rest/settings/' + mode, {device: disk_device}).done(on_complete).always(on_always).fail(on_error);
    };
    
backend.boot_extend = function(on_complete, on_error) {
        $.get('/rest/settings/boot_extend').done(on_complete).fail(on_error);
    };
