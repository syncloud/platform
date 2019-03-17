export function check_for_service_error(data, on_complete, on_error) {
    if (data.hasOwnProperty('success') && !data.success) {
        var xhr = {
            status: 200,
            responseJSON: data
        };
        on_error(xhr, {}, {});
    } else {
        on_complete();
    }
    
}

export function run_after_sam_is_complete(status_checker, timeout_func, on_complete, on_error) {
    run_after_job_is_complete(status_checker, timeout_func, on_complete, on_error, 'sam');
}

export function run_after_job_is_complete(status_checker, timeout_func, on_complete, on_error, job) {

    var recheck_function = function () { run_after_job_is_complete(status_checker, timeout_func, on_complete, on_error, job); };

    var recheck_timeout = 2000;
    status_checker(job,
        function (status) {
            if (status.is_running)
                timeout_func(recheck_function, recheck_timeout);
            else
                on_complete();
        },
        function (xhr, textStatus, errorThrown) {
            //Auth error means job is finished
            if (xhr.status == 401) {
                on_error(xhr, textStatus, errorThrown)
            } else {
                timeout_func(recheck_function, recheck_timeout);
            }
        }
    );
}

export function find_app(apps_data, app_id) {
    for (s=0; s < apps_data.length; s++) {
        var app_data = apps_data[s];
        if (app_data.app.id == app_id) {
            return app_data;
        }
    }
    return null;
}

export function get_value(values, name) {
    for (i=0; i < values.length; i++) {
        var value = values[i];
        if (value.name === name) {
            return value.value;
        }
    }
    return null;
}

export function job_status(job, on_complete, on_error) {
        $.get('/rest/settings/' + job + '_status').done(on_complete).fail(on_error);
    };

export function send_logs(include_support, on_always, on_error) {
        $.get('/rest/send_log',
          { include_support: include_support }
        ).always(on_always).fail(on_error);
    };