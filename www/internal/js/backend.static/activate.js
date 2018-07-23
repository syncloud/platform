var backend = {

    error_generic: {message: 'OMG!'},

    error_redirect_email: {parameters_messages: [{parameter: 'redirect_email', messages: ['redirect_email error']}]},
    error_redirect_password: {parameters_messages: [{parameter: 'redirect_password', messages: ['redirect_password error']}]},
    error_user_domain: {parameters_messages: [{parameter: 'user_domain', messages: ['user_domain error']}]},

    error_device_username: {parameters_messages: [{parameter: 'device_username', messages: ['device_username error']}]},
    error_device_password: {parameters_messages: [{parameter: 'device_password', messages: ['password error']}]},

    error_full_domain: {parameters_messages: [{parameter: 'full_domain', messages: ['full_domain error']}]},

    activate: function(parameters, on_always, on_done, on_error) {
        var that = this;
        setTimeout(function() {
            on_always();
            on_error({responseJSON: that.error_user_domain}, "error", {})
        }, 2000);
    },
    activate_custom_domain: function(parameters, on_always, on_done, on_error) {
        var that = this;
        setTimeout(function() {
            on_always();
            on_error({responseJSON: that.error_full_domain}, "error", {})
        }, 2000);
    },
    login: function() {
        window.location.href = "login.html";
    }
};