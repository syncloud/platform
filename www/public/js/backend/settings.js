var backend = {
    device_url: function(on_complete, on_error) {
        $.get('/rest/settings/device_domain').done(on_complete).fail(on_error);
    },

    send_logs: function(include_suppprt, on_always, on_error) {
        $.get('/rest/send_log',
          { include_suppprt: parameters.include_suppprt }
        ).always(on_always).fail(on_error);
    },

    reactivate: function() {
        window.location.href = (new URI()).port(81).directory("").filename("").query("");
    },

    get_versions: function(on_complete, on_always, on_error) {
        $.get('/rest/settings/versions').done(on_complete).always(on_always).fail(on_error);
    },

    check_versions: function(on_always, on_error) {
        $.get('/rest/check').always(on_always).fail(on_error);
    },

    platform_upgrade: function(on_complete, on_error) {
        $.get('/rest/settings/system_upgrade').done(on_complete).fail(on_error);
    },

    sam_upgrade: function(on_complete, on_error) {
        $.get('/rest/settings/sam_upgrade').done(on_complete).fail(on_error);
    },

    update_disks: function(on_complete, on_error) {
        $.get('/rest/settings/disks').done(on_complete).fail(on_error);
    },

    update_boot_disk: function(on_complete, on_always, on_error) {
        $.get('/rest/settings/boot_disk').done(on_complete).always(on_always).fail(on_error);
    },

    disk_action: function(disk_device, is_activate, on_always, on_error) {
        var mode = is_activate ? "disk_activate" : "disk_deactivate";
        $.get('/rest/settings/' + mode, {device: disk_device}).always(on_always).fail(on_error);
    },
    
    boot_extend: function(on_complete, on_error) {
        $.get('/rest/settings/boot_extend').done(on_complete).fail(on_error);
    },

    job_status: function (job, on_complete, on_error) {
        $.get('/rest/settings/' + job + '_status').done(on_complete).fail(on_error);
    }
};
