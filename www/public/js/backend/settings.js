function backend_device_url(on_complete) {
    $.get('/rest/settings/device_domain')
            .done(function (data) {
                display_device_url(data);
                on_complete();
            })
            .fail(onError);
}

function backend_update_disks(on_complete) {
    $.get('/rest/settings/disks')
            .done(function (data) {
                display_disks(data);
                on_complete();
            })
            .fail(onError);
}

function backend_disk_action(disk_device, is_activate, on_complete) {
    var mode = is_activate ? "disk_activate" : "disk_deactivate";
    $.get('/rest/settings/' + mode, {device: disk_device})
            .done(function () {
                backend_update_disks(on_complete);
            })
            .fail(onError);
}

function backend_send_logs(on_complete) {
    $.get('/rest/send_log').always(on_complete);
}

function backend_reactivate() {
    var internal_web = (new URI()).port(81).directory("").filename("").query("");
    window.location.href = internal_web;
}

function backend_check_access(on_complete) {
    $.get('/rest/settings/access')
            .done(function (data) {
                display_access(data.data);
            })
            .fail(onError)
            .always(on_complete);
}

function backend_external_access(state, on_complete) {
    $.get('/rest/settings/set_external_access?external_access=' + state)
            .done(function (data) {
                backend_check_access(on_complete);
            })
            .fail(onError);
}

function backend_protocol(new_protocol, on_complete) {
    $.get('/rest/settings/set_protocol?protocol=' + new_protocol)
            .done(function (data) {
                backend_check_access(on_complete);
            })
            .fail(onError);
}

function update_versions(on_complete) {
    $.get('/rest/settings/versions')
            .done(function (data) {
                display_versions(data);
            })
            .fail(onError)
            .always(function() {
            		typeof on_complete === 'function' && on_complete();
            });
}

function backend_check_versions(on_complete) {
    $.get('/rest/check')
            .always(function() {
                run_after_sam_is_complete(function() {
                        update_versions(on_complete);
                });
            });
}

function backend_platform_upgrade(on_complete) {
    $.get('/rest/settings/system_upgrade')
            .always(function() {
                run_after_sam_is_complete(function() {
                    update_versions(on_complete);
                });
            });
}

function backend_sam_upgrade(on_complete) {
    $.get('/rest/settings/sam_upgrade')
            .always(function() {
                run_after_sam_is_complete(function() {
                    update_versions(on_complete);
                });
            });
}