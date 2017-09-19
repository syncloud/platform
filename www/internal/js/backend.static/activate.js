var backend = {
//    activate: function(parameters) {
//        setTimeout(function() {
//            success_callbacks(parameters);
//        }, 2000);
//    }

    activate: function(parameters) {
        setTimeout(function() {
            parameters.always();
            parameters.fail(400, {message: parameters.values.domain_type})
        }, 2000);
    },
    login: function(name, password) {
        window.location.href = "login.html";
    }
};