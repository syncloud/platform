backend.installed_apps = function (on_complete, on_error) {
    $.get('/rest/installed_apps').done(on_complete).fail(on_error);
};