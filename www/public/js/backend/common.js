backend.job_status = function (job, on_complete, on_error) {
        $.get('/rest/settings/' + job + '_status').done(on_complete).fail(on_error);
    };