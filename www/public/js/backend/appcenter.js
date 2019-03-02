backend.available_apps = function (on_complete, on_error) {
    $.get('/rest/available_apps').done(on_complete).fail(on_error);
};