function onError(xhr, textStatus, errorThrown) {
    if (xhr.status === 401) {
        window.location.href = "/login.html";
    } else {
        window.location.href = "/error.html";
    }
}

function run_after_sam_is_complete(on_complete) {

    var recheck_function = function () { run_after_sam_is_complete(on_complete); }

    var recheck_timeout = 2000;
    $.get('/rest/settings/sam_status')
            .done(function(sam) {
                if (sam.is_running)
                    setTimeout(recheck_function, recheck_timeout);
                else
                    on_complete();
            })
            .fail(function() {
                setTimeout(recheck_function, recheck_timeout);
            })
}

function find_app(apps_data, app_id) {
    for (s=0; s < apps_data.length; s++) {
        var app_data = apps_data[s];
        if (app_data.app.id == app_id) {
            return app_data;
        }
    }
    return null;
}