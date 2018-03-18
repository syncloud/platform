var backend = {
    activate: function(parameters, on_always, on_done, on_error) {
        setTimeout(function() {
            on_always();
            on_error({responseJSON: {message: 'OMG!'}}, "error", {})
        }, 2000);
    },
    activate_custom_domain: function(parameters, on_always, on_done, on_error) {
        setTimeout(function() {
            on_always();
            on_error({responseJSON: {parameters_messages: [{parameter: 'device_username', messages: ['error']}]}}, "error", {})
        }, 2000);
    },
    login: function(name, password) {
        window.location.href = "login.html";
    }
};