function onError(xhr, textStatus, errorThrown) {
    if (xhr.status === 401) {
        window.location.href = "/login.html";
    } else if (xhr.status === 0) {
        console.log('user navigated away from the page');
    } else {
        window.location.href = "/error.html";
    }
}

function check_for_service_error(data, parameters, on_complete) {
    
    if (data.hasOwnProperty('success') && !data.success) {
        if (parameters.hasOwnProperty("fail")) {
            parameters.fail(200, data);
        }
    } else {
        on_complete();
    }
    
}

function run_after_sam_is_complete(on_complete, on_error) {
    run_after_job_is_complete(on_complete, on_error, 'sam');
}

function run_after_boot_extend_is_complete(on_complete) {
    run_after_job_is_complete(on_complete, on_error, 'boot_extend');
}

function run_after_job_is_complete(on_complete, on_error, job) {

    var recheck_function = function () { run_after_job_is_complete(on_complete, job); };

    var recheck_timeout = 2000;
    $.get('/rest/settings/' + job + '_status')
            .done(function(status) {
                if (status.is_running)
                    setTimeout(recheck_function, recheck_timeout);
                else
                    on_complete();
            })
            .fail(function(xhr, textStatus, errorThrown) {
                //Auth error means job is finished
                if (xhr.status == 401){
                    on_error(xhr, textStatus, errorThrown)
                } else {
                    setTimeout(recheck_function, recheck_timeout);
                }
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

function get_value(values, name) {
    for (i=0; i < values.length; i++) {
        var value = values[i];
        if (value.name === name) {
            return value.value;
        }
    }
    return null;
}