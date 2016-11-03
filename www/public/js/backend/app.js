var backend = {
    load_app: function(parameters) {
        var app_id = parameters.app_id;
        $.get( '/rest/app', {app_id: app_id})
            .done(function (data) {
                if (parameters.hasOwnProperty("done")) {
                    parameters.done(data);
                }
            })
            .fail(function (xhr, textStatus, errorThrown) {
                var error = null;
                if (xhr.hasOwnProperty('responseJSON')) {
                    var error = xhr.responseJSON;
                }
                if (parameters.hasOwnProperty("fail")) {
                    parameters.fail(xhr.status, error);
                }
            })
            .always(function() {
                if (parameters.hasOwnProperty("always")) {
                    parameters.always();
                }
            });
    },
    app_action: function(parameters) {
        var action = parameters.action;
        var app_id = parameters.app_id;
        $.get("/rest/"+action, {app_id: app_id})
            .done(function(data) {
                check_for_service_error(data, parameters, function() {
                    run_after_sam_is_complete(function () {
                        if (parameters.hasOwnProperty("done")) {
                            parameters.done(data);
                        }
                    });
                });
            })
            .fail(function (xhr, textStatus, errorThrown) {
                var error = null;
                if (xhr.hasOwnProperty('responseJSON')) {
                    var error = xhr.responseJSON;
                }
                if (parameters.hasOwnProperty("fail")) {
                    parameters.fail(xhr.status, error);
                }
            })
            .always(function() {
                run_after_sam_is_complete(function() {
                    if (parameters.hasOwnProperty("always")) {
                        parameters.always();
                    }
                });
            });
    }
};