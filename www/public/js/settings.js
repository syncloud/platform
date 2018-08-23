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