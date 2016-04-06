function external_access_status(accessBtn) {
    $.get('/rest/settings/external_access')
            .done(function (data) {
                accessBtn.text(data.external_access ? 'Enabled' : 'Disabled');
                accessBtn.data("external_access", data.external_access);
                accessBtn.prop('disabled', false);
            }).fail(onError);
}

function protocol_status(protocolBtn) {
    $.get('/rest/settings/protocol')
            .done(function (data) {

                protocolBtn.text(data.protocol);
                protocolBtn.data("protocol", data.protocol);
                protocolBtn.prop('disabled', false);

            }).fail(onError);
}

function disks_status() {
    $.get('/rest/settings/disks')
            .done(function (data) {
                var template = $("#disks-template").html();
                $("#disks").html(_.template(template)(data));
            }).fail(onError);
}

function check_system_version() {
    $.get('/rest/settings/version')
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