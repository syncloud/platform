backend.job_status = function (job, on_complete, on_error) {
    backend.test_timeout(function() { on_complete({success: true, is_running: false}) }, 2000);
};

backend.test_timeout = function(on_complete, timeout) {
    if (backend.async) {
        setTimeout(on_complete, timeout);
    } else {
        on_complete();
    }
};