var backend = {
    send_logs: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    }
}