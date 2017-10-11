function check_versions(on_complete, on_error) {

    backend.check_versions(function () {
        run_after_sam_is_complete(
            backend.job_status,
            setTimeout,
            function () {
                backend.get_versions(
                    on_complete,
                    on_error);
            }, on_error);
        }, on_error);
}

function platform_upgrade(on_complete, on_error) {

    backend.platform_upgrade(function (data) {
        check_for_service_error(data, function () {
            run_after_sam_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.get_versions(
                         on_complete,
                         on_error);
                 }, on_error);
        }, on_error);
    }, on_error);
    
}

function boot_extend(on_complete, on_error) {

    backend.boot_extend(function (data) {
        check_for_service_error(data, function () {
            run_after_boot_extend_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.update_boot_disk(
                        on_complete,
                        on_error);
                }, on_error);
        }, on_error);
    }, on_error);

}

function sam_upgrade(on_complete, on_error) {

    backend.sam_upgrade(function (data) {
        check_for_service_error(data, function () {
            run_after_sam_is_complete(
                backend.job_status,
                setTimeout,
                function () {
                    backend.get_versions(
                        on_complete,
                        on_error);
                }, on_error);
        }, on_error);
    }, on_error);
    
}

function disk_action(disk_device, is_activate, on_always, on_error) {
    on_complete = function(data) { check_for_service_error(data, function() {}, on_error); }
    backend.disk_action(disk_device, is_activate, on_complete, on_always, on_error);
}
