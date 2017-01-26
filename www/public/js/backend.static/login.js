var backend = {
    login: function(parameters) {
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    }

//    login: function(parameters) {
//        setTimeout(function() {
//            parameters.always();
//            parameters.fail(400, {parameters_messages: [{parameter: "name", messages: ["Login name can't be empty"]}]});
//        }, 2000);
//    }

}