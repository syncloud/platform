function backend_app_action(app_id, action, on_complete) {
    $.get("/rest/"+action, {app_id: app_id})
        .always(function() {
            run_after_sam_is_complete(function() {
                backend_update_app(app_id, on_complete);
            });
        });
}

function backend_update_app(app_id, on_complete) {
    $.get( '/rest/app', {app_id: app_id})
        .done( function(data) {
            ui_display_app(data);
        })
        .always(on_complete);
}