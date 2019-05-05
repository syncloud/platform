import URI from "urijs";
import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import { AppTemplate } from './app.templates.js'
import * as UiCommon from './ui/common.js'
import './ui/menu.js'

import * as Common from './common.js'

function load_app(app_id, on_complete, on_error) {
    $.get('/rest/app', {app_id: app_id}).done(on_complete).fail(on_error);
};

export function run_app_action(url, app_id, status_url, status_predicate, on_complete, on_error) {
    $.get(url, { app_id: app_id })
        .always((data) => {
            Common.check_for_service_error(data, () => {
                Common.run_after_job_is_complete(
                    setTimeout,
                    on_complete,
                    on_error,
                    status_url,
                    status_predicate);
            }, on_error)
        })
        .fail(on_error);
}

function register_btn_open_click() {
    $("#btn_open").off('click').on('click', function () {
        var btn = $(this);
        var app_url = btn.data('url');
        window.location.href = app_url;
    });
}

function register_btn_action_click(name, url) {
    const action = name.toLowerCase();

    $("#btn_" + action).off('click').on('click', function () {
         $('#app_action').val(action);
         $('#app_action_url').val(url);
         $('#confirm_caption').html(name);
         $('#app_action_confirmation').modal('show');
    });
}

function ui_display_app(data) {
		$("#block_app").html(_.template(AppTemplate)(data));
		const app_id = data.info.app.id;
		register_btn_open_click();
		register_btn_action_click('Install', '/rest/install');
		register_btn_action_click('Upgrade', '/rest/upgrade');
		register_btn_action_click('Remove', '/rest/remove');
	
 $("#btn_backup").off('click').on('click', function () {
         $('#backup_confirmation').modal('show');
    });

 $("#btn_backup_confirm").off('click').on('click', function () {
        var btn = $("#btn_backup");
        btn.button('loading');
       
        $.get('/rest/backup/create', {app: app_id})
         .always((data) => {
            Common.check_for_service_error(data, () => {
                Common.run_after_job_is_complete(
                    setTimeout,
                    () => {
                        btn.button('reset');
                        ui_load_app();
                    },
                    UiCommon.ui_display_error,
                    Common.JOB_STATUS_URL,
                    Common.JOB_STATUS_PREDICATE);
            }, UiCommon.ui_display_error)
         })
         .fail(UiCommon.ui_display_error);
  });

	 $("#btn_confirm").off('click').on('click', function () {
        var btn = $("#btn_" + $('#app_action').val());
        btn.button('loading');

        run_app_action(
            $('#app_action_url').val(),
            app_id,
            Common.INSTALLER_STATUS_URL,
            Common.DEFAULT_STATUS_PREDICATE,
            () => {
                btn.button('reset');
                ui_load_app();
            }, 
            UiCommon.ui_display_error);
  });
}

function ui_load_app() {
		var app_id = new URI().query(true)['app_id'];

		load_app(app_id, ui_display_app, UiCommon.ui_display_error);
}

$( document ).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };
    ui_load_app();
});
