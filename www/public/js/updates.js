import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import * as UiCommon from './ui/common.js'
import './ui/menu.js'

import * as Common from './common.js'

function get_versions(on_complete, on_error) {
    $.get('/rest/settings/versions').done(on_complete).fail(on_error);
};

export function check_versions(on_complete, on_error) {
    $.get('/rest/check')
        .always(function () {
            Common.run_after_job_is_complete(
                setTimeout,
                function () {
                    get_versions(
                        on_complete,
                        on_error);
                }, 
                on_error, 
                Common.INSTALLER_STATUS_URL,
                Common.DEFAULT_STATUS_PREDICATE);
            })
        .fail(on_error);
}

export function platform_upgrade(on_complete, on_error) {
    $.get('/rest/upgrade', { app_id: 'platform' })
        .done(function (data) {
                      Common.check_for_service_error(data, function () {
                          Common.run_after_job_is_complete(
                              setTimeout,
                              function () {
                                  get_versions(
                                       on_complete,
                                       on_error);
                               }, 
                               on_error,
                               Common.INSTALLER_STATUS_URL,
                               Common.DEFAULT_STATUS_PREDICATE);
                      }, on_error);
                  })
        .fail(on_error);

}

export function sam_upgrade(on_complete, on_error) {

    $.get('/rest/upgrade', { app_id: 'sam' })
        .done(function (data) {
                      Common.check_for_service_error(data, function () {
                          Common.run_after_job_is_complete(
                              setTimeout,
                              function () {
                                  get_versions(
                                      on_complete,
                                      on_error);
                              },
                              on_error,
                              Common.INSTALLER_STATUS_URL,
                              Common.DEFAULT_STATUS_PREDICATE);
                      }, on_error);
                  })
        .fail(on_error);

}

function ui_display_toggles() {
	$("[type='checkbox']").each(function() {
		$(this).bootstrapSwitch();
	});
}


function upgrade_buttons_enabled(is_enabled) {
		var btn_platform = $("#btn_platform_upgrade");
		var btn_sam = $("#btn_sam_upgrade");
		btn_platform.prop('disabled', !is_enabled);
		btn_sam.prop('disabled', !is_enabled);
}

function ui_display_versions(data) {

		var platform_data = Common.find_app(data.data, "platform");
		var sam_data = Common.find_app(data.data, "sam");

		$("#txt_platform_version").html(platform_data.installed_version);
		$("#txt_system_version_available").html(platform_data.current_version);

		if (platform_data.installed_version != platform_data.current_version) {
				$("#block_system_upgrade").show();
		} else {
				$("#block_system_upgrade").hide();
		}

		$("#txt_sam_version").html(sam_data.installed_version);
		$("#txt_sam_version_available").html(sam_data.current_version);

		if (sam_data.installed_version && sam_data.current_version && sam_data.installed_version != sam_data.current_version) {
				$("#block_sam_upgrade").show();
		} else {
				$("#block_sam_upgrade").hide();
		}
}

function ui_get_versions(on_always) {
		get_versions(ui_display_versions, on_always, UiCommon.ui_display_error);
}

function ui_check_versions() {
    var btn = $("#btn_check_updates");
    upgrade_buttons_enabled(false);
    btn.button('loading');
    
    check_versions(
        function (data) {

            ui_display_versions(data);
            btn.button('reset');
            upgrade_buttons_enabled(true);
        }, 
        function (a, b, c) {
            UiCommon.ui_display_error(a, b, c);
            btn.button('reset');
            upgrade_buttons_enabled(true);
        });
}

function ui_platform_upgrade() {
    var btn = $("#btn_platform_upgrade");
    btn.button('loading');

    platform_upgrade(
        function (data) {
            ui_display_versions(data);
            btn.button('reset');
        }, 
        function (a, b, c) {
            UiCommon.ui_display_error(a, b, c);
            btn.button('reset');
        });
 
}

function ui_sam_upgrade() {
    var btn = $("#btn_sam_upgrade");
    btn.button('loading');

    sam_upgrade(
        function (data) {
            ui_display_versions(data);
            btn.button('reset');
        },
        function (a, b, c) {
            UiCommon.ui_display_error(a, b, c);
            btn.button('reset');
        });
 
}

$(document).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };

    ui_display_toggles();

    $("#btn_check_updates").on('click', function () {
    		ui_check_versions();
    });

    $("#btn_platform_upgrade").on('click', function () {
    		ui_platform_upgrade();
    });

    $("#btn_sam_upgrade").on('click', function () {
    		ui_sam_upgrade();
    });

    ui_check_versions();
});
