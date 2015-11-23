function external_access_status(secureBtn, insecureBtn) {
    $.get('/server/rest/settings/external_access')
            .done(function (data) {

                secureBtn.text(data.mode == 'https' ? 'Enabled' : 'Disabled');
                secureBtn.data("enabled", data.mode == 'https');
                secureBtn.prop('disabled', false);

                insecureBtn.text(data.mode == 'http' ? 'Enabled' : 'Disabled');
                insecureBtn.data("enabled", data.mode == 'http');
                insecureBtn.prop('disabled', false);

            }).fail(onError);
}

function disks_status() {
    $.get('/server/rest/settings/disks')
            .done(function (data) {
                var template = $("#disks-template").html();
                $("#disks").html(_.template(template)(data));
            }).fail(onError);
}

function check_system_version() {
    $.get('/server/rest/settings/version')
            .done(function (data) {
                var btn = $("#system_upgrade_btn");
                var check_btn = $("#system_check_btn");
                btn.hide();
                check_btn.hide();
                var updates = '';

                if (data.installed_version != data.current_version) {
                    btn.show();
                    updates = ', updated version available: ' + data.current_version;
                    btn.text('Upgrade');
                    btn.prop('disabled', false);
                } else {
                    check_btn.show();
                    check_btn.text('Check');
                    check_btn.prop('disabled', false);
                }

                $("#system_version").html(data.installed_version + updates);

            }).fail(onError);
}

function external_access(btn, protocol, on_complete) {
    btn.prop('disabled', true);
    btn.text('Checking ...');
    var mode = btn.data("enabled") ? "disable" : "enable?mode=" + protocol;
    $.get('/server/rest/settings/external_access_' + mode)
            .done(function (data) {
                if (data.success)
                    $("#external_access_error").text('');
                else
                    $("#external_access_error").text(data.message);
                on_complete();
            }).fail(onError);
}