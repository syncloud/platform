function boot_extend(on_complete, on_error) {

    backend.boot_extend(function (data) {
        check_for_service_error(data, function () {
            run_after_job_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.update_boot_disk(
                        on_complete,
                        on_error);
                }, on_error, 'boot_extend');
        }, on_error);
    }, on_error);

}

function disk_action(disk_device, is_activate, on_always, on_error) {
    on_complete = function(data) { check_for_service_error(data, function() {}, on_error); }
    backend.disk_action(disk_device, is_activate, on_complete, on_always, on_error);
}


function disk_format(disk_device, on_complete, on_error) {
    backend.disk_format(disk_device, function (data) {
        check_for_service_error(data, function () {
            run_after_job_is_complete(
                backend.job_status,
                setTimeout,
                on_complete,
                on_error,
                'disk_format');
        }, on_error);
    }, on_error);
}
