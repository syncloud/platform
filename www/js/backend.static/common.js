function find_app(apps_data, app_id) {
    for (s=0; s < apps_data.length; s++) {
        var app_data = apps_data[s];
        if (app_data.app.id == app_id) {
            return app_data;
        }
    }
    return null;
}

function on_error(status, error) {
    window.location.href = "/error.html";
}

function success_callbacks(parameters, data) {
    if (parameters.hasOwnProperty("done")) {
        parameters.done(data);
    }
    if (parameters.hasOwnProperty("always")) {
        parameters.always();
    }
}