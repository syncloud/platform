var backend = {
    send_log: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    }
}