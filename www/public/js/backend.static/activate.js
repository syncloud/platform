var backend = {
    activate: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    }
}