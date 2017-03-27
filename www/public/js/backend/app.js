backend.load_app = function (app_id, on_complete, on_error) {
    $.get('/rest/app', {app_id: app_id}).done(on_complete).fail(on_error);
};

backend.app_action = function (app_id, action, on_always, on_error) {
    $.get("/rest/" + action, {app_id: app_id}).always(on_always).fail(on_error);
};