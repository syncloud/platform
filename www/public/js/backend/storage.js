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
    
backend.disk_format = function(disk_device, on_complete, on_error) {
        $.post('/rest/storage/disk_format', {device: disk_device}).done(on_complete).fail(on_error);
    };
 
