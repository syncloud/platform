function check_versions(un_complete, on_always, on_error) {

    backend.check_versions(function () {
                run_after_sam_is_complete(function () {
                            backend.get_versions(
                            un_complete, 
                            on_always, 
                            on_error);
                        },
                        on_error
                );
            },
            on_error);
}