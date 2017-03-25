function run_app_action(app_id, action, on_complete, on_error) {

    backend.app_action(app_id, action, function (data) {
        check_for_service_error(data, function () {
            run_after_sam_is_complete(
                on_complete,
                on_error);
        }, on_error)
    }, on_error);
}