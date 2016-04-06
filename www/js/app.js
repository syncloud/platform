function get_actions(info) {
    actions = [];
    if (info.installed_version) {
        actions.push('open');
        if (info.current_version != info.installed_version) {
            actions.push('upgrade');
        }
        actions.push('remove');
    } else {
        actions.push('install');
    }
    actions.push('check');

    return actions;
}

function progress_start() {
    $('#progress-bar').css('visibility', 'visible');
}

function progress_stop() {
    $('#progress-bar').css('visibility', 'hidden');
}

function refresh(app_id) {
    $.get( '/rest/app', {app_id: app_id})
            .done( function(data) {
                var template = $("#app-template").html();
                data.actions = get_actions(data.info);
                $("#app").html(_.template(template)(data));
                register_events(app_id);
                register_url(data.info.app.url, 'open');
            }).fail( onError );

}

function app_status(app_id) {
    refresh(app_id);
    progress_stop();
}

function run(action, app_id) {
    $('#errors_placeholder').empty();
    progress_start();
    $.get("/rest/" + action, {app_id: app_id})
            .fail(function(xhr, textStatus, errorThrown) {
                var template = $("#error_template").html();
                $('#errors_placeholder').html(_.template(template)(xhr.responseJSON));
            })
            .done(function() {
                run_after_sam_is_complete(function() {app_status(app_id)});
            });
}

function register_events(app_id) {
    register_event(app_id, 'install');
    register_event(app_id, 'remove');
    register_event(app_id, 'upgrade');
    register_event(app_id, 'check');
}

function register_event(app_id, action) {
    var btn = $("#btn" + action);
    btn.click(function (event) {
        btn.addClass('disabled');
        event.preventDefault();
        run(action, app_id);
    });
}

function register_url(url, action) {
    var btn = $("#btn" + action);
    btn.click(function (event) {
        event.preventDefault();
        window.location.href = url
    });
}