function available_apps(on_complete, on_error) {

    backend.available_apps(
         function (data) {
            check_for_service_error(
                data,
                function() {
                    on_complete(data);
                },
                on_error);
         },
         on_error
    );
}